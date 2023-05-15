package data

import (
	"github.com/awakari/core/api/grpc/writer"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var Msgs = writer.SubmitMessagesRequest{
	Msgs: []*pb.CloudEvent{
		{
			Id:          uuid.NewString(),
			Source:      "http://arxiv.org/abs/2305.06364",
			SpecVersion: "1.0",
			Type:        "com.github.awakari.producer-rss",
			Attributes: map[string]*pb.CloudEventAttributeValue{
				"summary": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "<p>We propose that the dark matter of our universe could be sterile neutrinos which reside within the twin sector of a mirror twin Higgs model. In our scenario, these particles are produced through a version of the Dodelson-Widrow mechanism that takes place entirely within the twin sector, yielding a dark matter candidate that is consistent with X-ray and gamma-ray line constraints. Furthermore, this scenario can naturally avoid the cosmological problems that are typically encountered in mirror twin Higgs models. In particular, if the sterile neutrinos in the Standard Model sector decay out of equilibrium, they can heat the Standard Model bath and reduce the contributions of the twin particles to $N_\\mathrm{eff}$. Such decays also reduce the effective temperature of the dark matter, thereby relaxing constraints from large-scale structure. The sterile neutrinos included in this model are compatible with the seesaw mechanism for generating Standard Model neutrino masses. </p> ",
					},
				},
				"tags": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "neutrino dark matter cosmology higgs standard model dodelson-widrow",
					},
				},
				"title": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Twin Sterile Neutrino Dark Matter. (arXiv:2305.06364v1 [hep-ph])",
					},
				},
			},
			Data: &pb.CloudEvent_TextData{
				TextData: "",
			},
		},
		{
			Id:          uuid.NewString(),
			Source:      "https://www.bbc.co.uk/news/business-65574826?at_medium=RSS&amp;at_campaign=KARANGA",
			SpecVersion: "1.0",
			Type:        "com.github.awakari.producer-rss",
			Attributes: map[string]*pb.CloudEventAttributeValue{
				"summary": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "NBCUniversal's former head of advertising is revealed as the new boss of the social network.",
					},
				},
				"time": {
					Attr: &pb.CloudEventAttributeValue_CeTimestamp{
						CeTimestamp: timestamppb.New(time.Date(2023, 05, 12, 16, 26, 45, 0, time.UTC)),
					},
				},
				"title": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Elon Musk names Linda Yaccarino new Twitter CEO",
					},
				},
			},
			Data: &pb.CloudEvent_TextData{
				TextData: "",
			},
		},
		{
			Id:          uuid.NewString(),
			Source:      "https://lenta.ru/news/2023/05/12/mjfox_depression/",
			SpecVersion: "1.0",
			Type:        "com.github.awakari.producer-rss",
			Attributes: map[string]*pb.CloudEventAttributeValue{
				"author": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Наталья Обрядина",
					},
				},
				"categories": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Забота о себе",
					},
				},
				"language": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "ru",
					},
				},
				"summary": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Американский актер Майкл Джей Фокс раскрыл свой способ не впасть в депрессию из-за прогрессирующей болезни Паркинсона. Он признался, что его состояние медленно, но верно ухудшается несмотря на лечение, однако он не теряет оптимизма. В этом звезде трилогии «Назад в будущее» помогает правило жить одним днем.",
					},
				},
				"time": {
					Attr: &pb.CloudEventAttributeValue_CeTimestamp{
						CeTimestamp: timestamppb.New(time.Date(2023, 05, 12, 17, 14, 18, 0, time.UTC)),
					},
				},
				"title": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Звезда «Назад в будущее» раскрыл способ избежать депрессии",
					},
				},
			},
			Data: &pb.CloudEvent_TextData{
				TextData: "",
			},
		},
		{
			Id:          uuid.NewString(),
			Source:      "https://yle.fi/a/74-20031356?origin=rss",
			SpecVersion: "1.0",
			Type:        "com.github.awakari.producer-rss",
			Attributes: map[string]*pb.CloudEventAttributeValue{
				"categories": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "musiikki Euroviisut Diivat popmusiikki",
					},
				},
				"imageurl": {
					Attr: &pb.CloudEventAttributeValue_CeUri{
						CeUri: "https://images.cdn.yle.fi/image/upload//w_205,h_115,q_70/39-1111768645bdce1cbadf.jpg",
					},
				},
				"language": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "fi",
					},
				},
				"summary": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Ruotsin edustaja on kääntänyt kelkkansa myös Suomen Käärijän suhteen. Loreen arvostaa Käärijän aitoutta. ",
					},
				},
				"time": {
					Attr: &pb.CloudEventAttributeValue_CeTimestamp{
						CeTimestamp: timestamppb.New(time.Date(2023, 05, 12, 16, 06, 04, 0, time.UTC)),
					},
				},
				"title": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Euroviisujen ennakkosuosikki Loreen on rakastettu ja vihattu diiva, mutta Käärijälle häneltä heruu kehuja",
					},
				},
			},
			Data: &pb.CloudEvent_TextData{
				TextData: "",
			},
		},
		{
			Id:          uuid.NewString(),
			Source:      "github.com/awakari/tgbot",
			SpecVersion: "1.0",
			Type:        "com.github.awakari.tgbot",
			Attributes: map[string]*pb.CloudEventAttributeValue{
				"author": {
					Attr: &pb.CloudEventAttributeValue_CeString{
						CeString: "Edna",
					},
				},
				"time": {
					Attr: &pb.CloudEventAttributeValue_CeTimestamp{
						CeTimestamp: timestamppb.New(time.Date(2023, 05, 12, 18, 38, 45, 0, time.UTC)),
					},
				},
			},
			Data: &pb.CloudEvent_TextData{
				TextData: "hi there",
			},
		},
	},
}
