package data

import "github.com/awakari/core/api/grpc/subscriptions"

var Subs = []subscriptions.CreateRequest{
	{
		Md: &subscriptions.Metadata{
			Description: "disabled",
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Ktc{
				Ktc: &subscriptions.KiwiTreeConditionInput{
					Key:     "author",
					Pattern: "Edna",
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "exact complete value match for a key",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Ktc{
				Ktc: &subscriptions.KiwiTreeConditionInput{
					Key:     "author",
					Pattern: "Edna",
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "partial exact match",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Ktc{
				Ktc: &subscriptions.KiwiTreeConditionInput{
					Key:     "tags",
					Pattern: "neutrino",
					Partial: true,
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "basic group condition with \"and\" logic and partial sub-conditions",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Gc{
				Gc: &subscriptions.GroupConditionInput{
					Logic: 0,
					Group: []*subscriptions.ConditionInput{
						{
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "title",
									Pattern: "Elon",
									Partial: true,
								},
							},
						},
						{
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "title",
									Pattern: "Musk",
									Partial: true,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "basic group condition with \"or\" logic",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Gc{
				Gc: &subscriptions.GroupConditionInput{
					Logic: 1,
					Group: []*subscriptions.ConditionInput{
						{
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "language",
									Pattern: "fi",
								},
							},
						},
						{
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "language",
									Pattern: "ru",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "basic group condition with \"and\" logic and a negative sub-condition",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Gc{
				Gc: &subscriptions.GroupConditionInput{
					Logic: 0,
					Group: []*subscriptions.ConditionInput{
						{
							Not: true,
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "type",
									Pattern: "com.github.awakari.tgbot",
								},
							},
						},
						{
							Cond: &subscriptions.ConditionInput_Ktc{
								Ktc: &subscriptions.KiwiTreeConditionInput{
									Key:     "summary",
									Pattern: "of",
									Partial: true,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "single symbol wildcard",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Ktc{
				Ktc: &subscriptions.KiwiTreeConditionInput{
					Key:     "title",
					Pattern: "?eutrino",
					Partial: true,
				},
			},
		},
	},
	{
		Md: &subscriptions.Metadata{
			Description: "multiple symbol wildcard",
			Enabled:     true,
		},
		Cond: &subscriptions.ConditionInput{
			Cond: &subscriptions.ConditionInput_Ktc{
				Ktc: &subscriptions.KiwiTreeConditionInput{
					Key:     "time",
					Pattern: "2023-05-12T16*",
				},
			},
		},
	},
}
