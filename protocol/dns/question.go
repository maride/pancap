package dns

import (
	"fmt"
	"github.com/maride/pancap/common"
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
func (p *Protocol) processDNSQuestion(questions []layers.DNSQuestion) {
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
		p.processType(questionType, question.Type)

		// Append full domain and base domain
		questionDomains = common.AppendIfUnique(name, questionDomains)

		// Check if we need to add the base name to the private list
		_, icannManaged := publicsuffix.PublicSuffix(name)
		if icannManaged {
			// TLD is managed by ICANN, add to the base list
			questionBaseDomains = common.AppendIfUnique(basename, questionBaseDomains)
		} else {
			// it's not managed by ICANN, so it's private - add it to the private list
			questionPrivateDomains = common.AppendIfUnique(name, questionPrivateDomains)
		}
	}
}

// Generates a summary of all DNS questions
func (p *Protocol) generateDNSQuestionSummary() string {
	summary := ""

	// Overall question stats
	summary = fmt.Sprintf("%s%d DNS questions in total\n", summary, numQuestions)
	summary = fmt.Sprintf("%s%s records\n", summary, p.generateDNSTypeSummary(questionType))
	summary = fmt.Sprintf("%s%d unique domains of %d base domains, of which are %d private (non-ICANN) TLDs.\n", summary, len(questionDomains), len(questionBaseDomains), len(questionPrivateDomains))

	// Output base domains asked for
	if len(questionBaseDomains) > 0 {
		summary = fmt.Sprintf("%sAsked for these base domains:\n%s", summary, common.GenerateTree(questionBaseDomains))
	}

	// Output private domains
	if len(questionPrivateDomains) > 0 {
		summary = fmt.Sprintf("%sAsked for these private (non-ICANN managed) domains:\n%s", summary, common.GenerateTree(questionPrivateDomains))
	}

	// And return summary
	return summary
}