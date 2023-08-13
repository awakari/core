package data

import (
	"github.com/awakari/client-sdk-go/model/subscription"
	"github.com/awakari/client-sdk-go/model/subscription/condition"
	"time"
)

var Subs = []subscription.Data{
	{
		Description: "disabled",
		Condition: condition.
			NewBuilder().
			AttributeKey("author").
			AnyOfWords("Edna").
			BuildTextCondition(),
	},
	{
		Description: "exact complete value match for a key",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			AttributeKey("author").
			TextEquals("Edna").
			BuildTextCondition(),
	},
	{
		Description: "partial exact match",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			AttributeKey("tags").
			AnyOfWords("neutrino").
			BuildTextCondition(),
	},
	{
		Description: "basic group condition with \"and\" logic and partial sub-conditions",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			All([]condition.Condition{
				condition.
					NewBuilder().
					AttributeKey("title").
					AnyOfWords("Elon").
					BuildTextCondition(),
				condition.
					NewBuilder().
					AttributeKey("title").
					AnyOfWords("Musk").
					BuildTextCondition(),
			}).
			BuildGroupCondition(),
	},
	{
		Description: "basic group condition with \"or\" logic",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			Any([]condition.Condition{
				condition.
					NewBuilder().
					AttributeKey("language").
					AnyOfWords("fi").
					BuildTextCondition(),
				condition.
					NewBuilder().
					AttributeKey("language").
					AnyOfWords("ru").
					BuildTextCondition(),
			}).
			BuildGroupCondition(),
	},
	{
		Description: "basic group condition with \"and\" logic and a negative sub-condition",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			All([]condition.Condition{
				condition.
					NewBuilder().
					Not().
					AttributeKey("type").
					TextEquals("com.github.awakari.tgbot").
					BuildTextCondition(),
				condition.
					NewBuilder().
					AttributeKey("summary").
					AnyOfWords("propose").
					BuildTextCondition(),
			}).
			BuildGroupCondition(),
	},
	{
		Description: "before a certain time",
		Enabled:     true,
		Condition: condition.
			NewBuilder().
			AttributeKey("time").
			LessThan(float64(time.Date(2023, 05, 12, 17, 0, 0, 0, time.UTC).UnixMilli())).
			BuildNumberCondition(),
	},
}
