package dns

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/publicsuffix"
	"log"
)

var (
	numAnswers int
	answerDomains []string
	answerBaseDomains []string
	answerPrivateDomains []string
	answerType = make(map[layers.DNSType]int)
	answerPublicIPv4 []string
	answerPrivateIPv4 []string
)

// Called on every DNS packet to process response(s)
func processDNSAnswer(answers []layers.DNSResourceRecord) {
	for _, answer := range answers {
		// Raise stats
		numAnswers++

		// Add answer to answers array
		name := string(answer.Name)
		basename, basenameErr := publicsuffix.EffectiveTLDPlusOne(name)

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
		_, icannManaged := publicsuffix.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			answerBaseDomains = appendIfUnique(basename, answerBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			answerPrivateDomains = appendIfUnique(name, answerPrivateDomains)
		}

		// Check if we got an A record answer
		if answer.Type == layers.DNSTypeA {
			// A record, check IP for being private
			if ipIsPrivate(answer.IP) {
				answerPrivateIPv4 = append(answerPrivateIPv4, answer.IP.String())
			} else {
				answerPublicIPv4 = append(answerPublicIPv4, answer.IP.String())
			}
		}
	}
}

// Prints a summary of all DNS answers
func printDNSAnswerSummary() {
	// Overall question stats
	fmt.Printf("%d DNS answers in total\n", numAnswers)
	fmt.Printf("%s records\n", generateDNSTypeSummary(answerType))
	fmt.Printf("%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.\n", len(answerDomains), len(answerBaseDomains), len(answerPrivateDomains))

	// Output base domains answered with
	if len(answerBaseDomains) > 0 {
		fmt.Println("Answered with these base domains:")
		printTree(answerBaseDomains)
	}

	// Output private domains
	if len(answerPrivateDomains) > 0 {
		fmt.Println("Answered with these private (non-ICANN managed) domains:")
		printTree(answerPrivateDomains)
	}

	// Check for public and private IPs
	fmt.Printf("Answered with %d public IP addresses and %d private IP addresses\n", len(answerPublicIPv4), len(answerPrivateIPv4))
	if len(answerPrivateIPv4) > 0 {
		fmt.Println("Private IP addresses in answer:")
		printTree(answerPrivateIPv4)
	}
}
