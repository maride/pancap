package dns

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/publicsuffix"
	"log"
)

var (
	numQuestions int
	questionDomains []string
	questionBaseDomains []string
	questionPrivateDomains []string
	questionType = make(map[layers.DNSType]int)
)

// Called on every DNS packet to process questions
func processDNSQuestion(questions []layers.DNSQuestion) {
	// Iterate over all questions
	for _, question := range questions {
		// Raise stats
		numQuestions++

		// Add question to questions array
		name := string(question.Name)
		basename, basenameErr := publicsuffix.EffectiveTLDPlusOne(name)

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
		_, icannManaged := publicsuffix.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			questionBaseDomains = appendIfUnique(basename, questionBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			questionPrivateDomains = appendIfUnique(name, questionPrivateDomains)
		}
	}
}

// Prints a summary of all DNS questions
func printDNSQuestionSummary() {
	// Overall question stats
	fmt.Printf("%d DNS questions in total\n", numQuestions)
	fmt.Printf("%s records\n", generateDNSTypeSummary(questionType))
	fmt.Printf("%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.\n", len(questionDomains), len(questionBaseDomains), len(questionPrivateDomains))

	// Output base domains asked for
	if len(questionBaseDomains) > 0 {
		fmt.Println("Asked for these base domains:")
		printTree(questionBaseDomains)
	}

	// Output private domains
	if len(questionPrivateDomains) > 0 {
		fmt.Println("Asked for these private (non-ICANN managed) domains:")
		printTree(questionPrivateDomains)
	}
}