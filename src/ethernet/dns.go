package ethernet

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	publicsuffix2 "golang.org/x/net/publicsuffix"
	"log"
)

var (
	numQuestions int
	questionDomains []string
	questionBaseDomains []string
	questionPrivateDomains []string
	questionType = make(map[layers.DNSType]int)

	numAnswers int
	answerDomains []string
	answerBaseDomains []string
	answerPrivateDomains []string
	answerType = make(map[layers.DNSType]int)
)

// Called on every DNS packet
func processDNSPacket(packet gopacket.Packet) error {
	var dnspacket layers.DNS

	// Decode raw packet into DNS
	decodeErr := dnspacket.DecodeFromBytes(packet.ApplicationLayer().LayerContents(), gopacket.NilDecodeFeedback)
	if decodeErr != nil {
		// Encountered an error during decoding, most likely a broken packet
		return decodeErr
	}

	// Further process the packet
	processDNSQuestion(dnspacket.Questions)
	processDNSAnswer(dnspacket.Answers)

	// No error encountered, return clean
	return nil
}

// Called on every DNS packet to process questions
func processDNSQuestion(questions []layers.DNSQuestion) {
	// Iterate over all questions
	for _, question := range questions {
		// Raise stats
		numQuestions++

		// Add question to questions array
		name := string(question.Name)
		basename, basenameErr := publicsuffix2.EffectiveTLDPlusOne(name)

		if basenameErr != nil {
			// Encountered error while checking for the basename
			log.Printf("Encountered error while checking '%s' domain for its basename: %s", name, basenameErr.Error())
			continue
		}

		// Process type questions
		processType(questionType, question.Type)

		// Append full domain and base domain
		questionDomains = appendIfUnique(name, questionDomains)

		// Check if we need to add the base name to the private list
		_, icannManaged := publicsuffix2.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			questionBaseDomains = appendIfUnique(basename, questionBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			questionPrivateDomains = appendIfUnique(name, questionPrivateDomains)
		}
	}
}

// Called on every DNS packet to process response(s)
func processDNSAnswer(answers []layers.DNSResourceRecord) {
	for _, answer := range answers {
		// Raise stats
		numAnswers++

		// Add answer to answers array
		name := string(answer.Name)
		basename, basenameErr := publicsuffix2.EffectiveTLDPlusOne(name)

		if basenameErr != nil {
			// Encountered error while checking for the basename
			log.Printf("Encountered error while checking '%s' domain for its basename: %s", name, basenameErr.Error())
			continue
		}

		// Process type answers
		processType(answerType, answer.Type)

		// Append full domain and base domain
		answerDomains = appendIfUnique(name, answerDomains)

		// Check if we need to add the base name to the private list
		_, icannManaged := publicsuffix2.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			answerBaseDomains = appendIfUnique(basename, answerBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			answerPrivateDomains = appendIfUnique(name, answerPrivateDomains)
		}
	}
}

// Processes the given dnstype and raises its stats in the given array
func processType(typearr map[layers.DNSType]int, dnstype layers.DNSType) {
	typearr[dnstype]++
}

// Appends the appendee to the array if it does not contain appendee yet
func appendIfUnique(appendee string, array []string) []string {
	// Iterate over all elements and check values
	for _, elem := range array {
		if elem == appendee {
			// ... found. Stop here
			return array
		}
	}

	// None found, append
	return append(array, appendee)
}

// Print a summary after all DNS packets were processed
func printDNSSummary() {
	printDNSQuestionSummary()
	printDNSAnswerSummary()
}

// Prints a summary of all DNS questions
func printDNSQuestionSummary() {
	// Overall question stats
	log.Printf("%d DNS questions in total", numQuestions)
	log.Printf("%s records", generateDNSTypeSummary(questionType))
	log.Printf("%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.", len(questionDomains), len(questionBaseDomains), len(questionPrivateDomains))

	// Output base domains asked for
	if len(questionBaseDomains) > 0 {
		log.Println("Asked for these base domains:")
		printTree(questionBaseDomains)
	}

	// Output private domains
	if len(questionPrivateDomains) > 0 {
		log.Println("Asked for these private (non-ICANN managed) domains:")
		printTree(questionPrivateDomains)
	}
}

// Prints a summary of all DNS answers
func printDNSAnswerSummary() {
	// Overall question stats
	log.Printf("%d DNS answers in total", numAnswers)
	log.Printf("%s records", generateDNSTypeSummary(answerType))
	log.Printf("%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.", len(answerDomains), len(answerBaseDomains), len(answerPrivateDomains))

	// Output base domains answered with
	if len(answerBaseDomains) > 0 {
		log.Println("Answered with these base domains:")
		printTree(answerBaseDomains)
	}

	// Output private domains
	if len(answerPrivateDomains) > 0 {
		log.Println("Answered with these private (non-ICANN managed) domains:")
		printTree(answerPrivateDomains)
	}
}

// Prints each element, along with a small ASCII tree
func printTree(strarr []string) {
	// iterate over each element
	for iter, elem := range strarr {
		// check if we got the last element
		if iter < len(strarr) - 1 {
			log.Printf("|- %s", elem)
		} else {
			log.Printf("'- %s\n\n", elem)
		}
	}
}

// Generates a summary string for DNS types in the given array
func generateDNSTypeSummary(typearr map[layers.DNSType]int) string {
	var answerarr []string

	// Iterate over all possible DNS types
	for iter, typeelem := range typearr {
		// Read amount of type hits for this type
		answerarr = append(answerarr, fmt.Sprintf("%d %s", typeelem, iter.String()))
	}

	// Check if we even processed a single type
	if len(answerarr) == 0 {
		// we didn't, strange.
		return "(no types encountered)"
	}

	// now, glue all array elements together
	answerstr := ""
	for iter, elem := range answerarr {
		// Check if we need to apply to proper sentence rules
		if iter == 0 {
			// We don't need to append yet
			answerstr = elem
		} else if iter == len(answerarr) - 1 {
			// Last element, use "and" instead of a comma
			answerstr = fmt.Sprintf("%s and %s", answerstr, elem)
		} else {
			// Some entry, just add it with a comma
			answerstr = fmt.Sprintf("%s, %s", answerstr, elem)
		}
	}

	return answerstr
}
