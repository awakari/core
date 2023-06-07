package data

import (
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
)

var Subs = []subscription.Data{
	{
		Metadata: subscription.Metadata{
			Description: "disabled",
		},
		Condition: condition.
			NewBuilder().
			MatchAttrKey("author").
			MatchAttrValuePattern("Edna").
			BuildKiwiTreeCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "exact complete value match for a key",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			MatchAttrKey("author").
			MatchAttrValuePattern("Edna").
			BuildKiwiTreeCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "partial exact match",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			MatchAttrValuePartial().
			MatchAttrKey("tags").
			MatchAttrValuePattern("neutrino").
			BuildKiwiTreeCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "basic group condition with \"and\" logic and partial sub-conditions",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicAnd).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						MatchAttrKey("title").
						MatchAttrValuePattern("Elon").
						MatchAttrValuePartial().
						BuildKiwiTreeCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("title").
						MatchAttrValuePattern("Musk").
						MatchAttrValuePartial().
						BuildKiwiTreeCondition(),
				},
			).
			BuildGroupCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "basic group condition with \"or\" logic",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicOr).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						MatchAttrKey("language").
						MatchAttrValuePattern("fi").
						BuildKiwiTreeCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("language").
						MatchAttrValuePattern("ru").
						BuildKiwiTreeCondition(),
				},
			).
			BuildGroupCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "basic group condition with \"and\" logic and a negative sub-condition",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			GroupLogic(condition.GroupLogicAnd).
			GroupChildren(
				[]condition.Condition{
					condition.
						NewBuilder().
						Negation().
						MatchAttrKey("type").
						MatchAttrValuePattern("com.github.awakari.tgbot").
						BuildKiwiTreeCondition(),
					condition.
						NewBuilder().
						MatchAttrKey("summary").
						MatchAttrValuePattern("of").
						MatchAttrValuePartial().
						BuildKiwiTreeCondition(),
				},
			).
			BuildGroupCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "single symbol wildcard",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			MatchAttrValuePartial().
			MatchAttrKey("title").
			MatchAttrValuePattern("?eutrino").
			BuildKiwiTreeCondition(),
	},
	{
		Metadata: subscription.Metadata{
			Description: "multiple symbol wildcard",
			Enabled:     true,
		},
		Condition: condition.
			NewBuilder().
			MatchAttrKey("time").
			MatchAttrValuePattern("2023-05-12T16*").
			BuildKiwiTreeCondition(),
	},
}
