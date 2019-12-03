package dns

import (
	"fmt"
	"git.darknebu.la/maride/pancap/common"
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
		answerDomains = common.AppendIfUnique(name, answerDomains)

		// Check if we need to add the base name to the private list
		_, icannManaged := publicsuffix.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			answerBaseDomains = common.AppendIfUnique(basename, answerBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			answerPrivateDomains = common.AppendIfUnique(name, answerPrivateDomains)
		}

		// Check if we got an A record answer
		if answer.Type == layers.DNSTypeA {
			// A record, check IP for being private
			if ipIsPrivate(answer.IP) {
				answerPrivateIPv4 = common.AppendIfUnique(answer.IP.String(), answerPrivateIPv4)
			} else {
				answerPublicIPv4 = common.AppendIfUnique(answer.IP.String(), answerPublicIPv4)
			}
		}
	}
}

// Generates a summary of all DNS answers
func generateDNSAnswerSummary() string {
	summary := ""

	// Overall question stats
	summary = fmt.Sprintf("%s%d DNS answers in total\n", summary, numAnswers)
	summary = fmt.Sprintf("%s%s records\n", summary, generateDNSTypeSummary(answerType))
	summary = fmt.Sprintf("%s%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.\n", summary, len(answerDomains), len(answerBaseDomains), len(answerPrivateDomains))

	// Output base domains answered with
	if len(answerBaseDomains) > 0 {
		summary = fmt.Sprintf("Answered with these base domains:\n%s", common.GenerateTree(answerBaseDomains))
	}

	// Output private domains
	if len(answerPrivateDomains) > 0 {
		summary = fmt.Sprintf("%sAnswered with these private (non-ICANN managed) domains:\n%s", summary, common.GenerateTree(answerPrivateDomains))
	}

	// Check for public and private IPs
	summary = fmt.Sprintf("%sAnswered with %d public IP addresses and %d private IP addresses\n", summary, len(answerPublicIPv4), len(answerPrivateIPv4))
	if len(answerPrivateIPv4) > 0 {
		summary = fmt.Sprintf("%sPrivate IP addresses in answer:\n%s", summary, common.GenerateTree(answerPrivateIPv4))
	}

	// Return summary
	return summary
}
