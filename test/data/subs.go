package data

import (
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
)

var Subs = []subscription.Data{
	{
		Description: "disabled",
		Condition: condition.
			NewBuilder().
			MatchAttrKey("author").
			MatchText("Edna").
			BuildTextCondition(),
	},
	{
		Description: "exact complete value match for a key",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			MatchAttrKey("author").
			MatchText("Edna").
			BuildTextCondition(),
	},
	{
		Description: "partial exact match",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			MatchAttrKey("tags").
			MatchText("neutrino").
			BuildTextCondition(),
	},
	{
		Description: "basic group condition with \"and\" logic and partial sub-conditions",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicAnd).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						MatchAttrKey("title").
						MatchText("Elon").
						BuildTextCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("title").
						MatchText("Musk").
						BuildTextCondition(),
				},
			).
			BuildGroupCondition(),
	},
	{
		Description: "basic group condition with \"or\" logic",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicOr).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						MatchAttrKey("language").
						MatchText("fi").
						BuildTextCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("language").
						MatchText("ru").
						BuildTextCondition(),
				},
			).
			BuildGroupCondition(),
	},
	{
		Description: "basic group condition with \"and\" logic and a negative sub-condition",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicAnd).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						Negation().
						MatchAttrKey("type").
						MatchText("com.github.awakari.tgbot").
						MatchExact().
						BuildTextCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("summary").
						MatchText("propose").
						BuildTextCondition(),
				},
			).
			BuildGroupCondition(),
	},
}
