// Package mirea provides registry definitions for MIREA university sources.
package mirea

import (
	"github.com/trueegorletov/analabit/core/source"
	"github.com/trueegorletov/analabit/core/source/mirea"
)

var Varsity = source.VarsityDefinition{
	Code:           "mirea",
	Name:           "МИРЭА",
	HeadingSources: sourcesList(),
}

func sourcesList() []source.HeadingSource {
	return []source.HeadingSource{
		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098830464818486"},
			BVIListIDs:            []string{"1829098830460624182"},
			TargetQuotaListIDs:    []string{"1829098830463769910"},
			DedicatedQuotaListIDs: []string{"1829098830462721334"},
			SpecialQuotaListIDs:   []string{"1829098830461672758"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098644022762806"},
			BVIListIDs:            []string{"1829098644019617078"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098644021714230"},
			SpecialQuotaListIDs:   []string{"1829098644020665654"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098280654478646"},
			BVIListIDs:            []string{"1829098280651332918"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098280653430070"},
			SpecialQuotaListIDs:   []string{"1829098280652381494"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098359733886262"},
			BVIListIDs:            []string{"1829098359728643382"},
			TargetQuotaListIDs:    []string{"1829098359731789110", "1829098359732837686"},
			DedicatedQuotaListIDs: []string{"1829098359730740534"},
			SpecialQuotaListIDs:   []string{"1829098359729691958"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098713741532470"},
			BVIListIDs:            []string{"1829098713736289590"},
			TargetQuotaListIDs:    []string{"1829098713739435318", "1829098713740483894"},
			DedicatedQuotaListIDs: []string{"1829098713738386742"},
			SpecialQuotaListIDs:   []string{"1829098713737338166"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098622271102262"},
			BVIListIDs:            []string{"1829098622258519350"},
			TargetQuotaListIDs:    []string{"1829098622262713654", "1829098622263762230", "1829098622264810806", "1829098622265859382", "1829098622266907958", "1829098622267956534", "1829098622269005110", "1829098622270053686", "1829185534518369590"},
			DedicatedQuotaListIDs: []string{"1829098622260616502"},
			SpecialQuotaListIDs:   []string{"1829098622259567926"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098428193316150"},
			BVIListIDs:            []string{"1829098428187024694"},
			TargetQuotaListIDs:    []string{"1829098428191218998", "1829098428192267574", "1829185420161719606"},
			DedicatedQuotaListIDs: []string{"1829098428189121846"},
			SpecialQuotaListIDs:   []string{"1829098428188073270"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098290594979126"},
			BVIListIDs:            []string{"1829098290589736246"},
			TargetQuotaListIDs:    []string{"1829098290592881974", "1829098290593930550"},
			DedicatedQuotaListIDs: []string{"1829098290591833398"},
			SpecialQuotaListIDs:   []string{"1829098290590784822"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098929896037686"},
			BVIListIDs:            []string{"1829098929890794806"},
			TargetQuotaListIDs:    []string{"1829098929894989110", "1829185975260028214"},
			DedicatedQuotaListIDs: []string{"1829098929892891958"},
			SpecialQuotaListIDs:   []string{"1829098929891843382"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098349764025654"},
			BVIListIDs:            []string{"1829098349759831350"},
			TargetQuotaListIDs:    []string{"1829184047658573110"},
			DedicatedQuotaListIDs: []string{"1829098349761928502"},
			SpecialQuotaListIDs:   []string{"1829098349760879926"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829101057973689654"},
			BVIListIDs:            []string{"1829101057972641078"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{},
			SpecialQuotaListIDs:   []string{},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099154690809142"},
			BVIListIDs:            []string{"1829099154687663414"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099154689760566"},
			SpecialQuotaListIDs:   []string{"1829099154688711990"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098982235708726"},
			BVIListIDs:            []string{"1829098982230465846"},
			TargetQuotaListIDs:    []string{"1829098982233611574", "1829098982234660150"},
			DedicatedQuotaListIDs: []string{"1829098982232562998"},
			SpecialQuotaListIDs:   []string{"1829098982231514422"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098939903647030"},
			BVIListIDs:            []string{"1829098939899452726"},
			TargetQuotaListIDs:    []string{"1829184162223889718"},
			DedicatedQuotaListIDs: []string{"1829098939901549878"},
			SpecialQuotaListIDs:   []string{"1829098939900501302"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098991581666614"},
			BVIListIDs:            []string{"1829098991577472310"},
			TargetQuotaListIDs:    []string{"1829185985294900534"},
			DedicatedQuotaListIDs: []string{"1829098991579569462"},
			SpecialQuotaListIDs:   []string{"1829098991578520886"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098767469518134"},
			BVIListIDs:            []string{"1829098767464275254"},
			TargetQuotaListIDs:    []string{"1829098767467420982", "1829098767468469558"},
			DedicatedQuotaListIDs: []string{"1829098767466372406"},
			SpecialQuotaListIDs:   []string{"1829098767465323830"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098778444963126"},
			BVIListIDs:            []string{"1829098778439720246"},
			TargetQuotaListIDs:    []string{"1829098778442865974", "1829098778443914550"},
			DedicatedQuotaListIDs: []string{"1829098778441817398"},
			SpecialQuotaListIDs:   []string{"1829098778440768822"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098545798454582"},
			BVIListIDs:            []string{"1829098545794260278"},
			TargetQuotaListIDs:    []string{"1829098545797406006"},
			DedicatedQuotaListIDs: []string{"1829098545796357430"},
			SpecialQuotaListIDs:   []string{"1829098545795308854"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098369850547510"},
			BVIListIDs:            []string{"1829098369845304630"},
			TargetQuotaListIDs:    []string{"1829098369848450358", "1829098369849498934"},
			DedicatedQuotaListIDs: []string{"1829098369847401782"},
			SpecialQuotaListIDs:   []string{"1829098369846353206"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098723858193718"},
			BVIListIDs:            []string{"1829098723848756534"},
			TargetQuotaListIDs:    []string{"1829098723852950838", "1829098723853999414", "1829098723855047990", "1829098723856096566", "1829098723857145142", "1829185800051367222"},
			DedicatedQuotaListIDs: []string{"1829098723850853686"},
			SpecialQuotaListIDs:   []string{"1829098723849805110"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098511654722870"},
			BVIListIDs:            []string{"1829098511648431414"},
			TargetQuotaListIDs:    []string{"1829098511651577142", "1829098511652625718", "1829098511653674294"},
			DedicatedQuotaListIDs: []string{"1829098511650528566"},
			SpecialQuotaListIDs:   []string{"1829098511649479990"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098557629537590"},
			BVIListIDs:            []string{"1829098557623246134"},
			TargetQuotaListIDs:    []string{"1829098557626391862", "1829098557627440438", "1829098557628489014"},
			DedicatedQuotaListIDs: []string{"1829098557625343286"},
			SpecialQuotaListIDs:   []string{"1829098557624294710"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098567545920822"},
			BVIListIDs:            []string{"1829098567540677942"},
			TargetQuotaListIDs:    []string{"1829098567543823670", "1829098567544872246"},
			DedicatedQuotaListIDs: []string{"1829098567542775094"},
			SpecialQuotaListIDs:   []string{"1829098567541726518"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099079887494454"},
			BVIListIDs:            []string{"1829099079884348726"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099079886445878"},
			SpecialQuotaListIDs:   []string{"1829099079885397302"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098380367764790"},
			BVIListIDs:            []string{"1829098380362521910"},
			TargetQuotaListIDs:    []string{"1829098380365667638", "1829098380366716214"},
			DedicatedQuotaListIDs: []string{"1829098380364619062"},
			SpecialQuotaListIDs:   []string{"1829098380363570486"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098300710591798"},
			BVIListIDs:            []string{"1829098300705348918"},
			TargetQuotaListIDs:    []string{"1829098300708494646", "1829098300709543222"},
			DedicatedQuotaListIDs: []string{"1829098300707446070"},
			SpecialQuotaListIDs:   []string{"1829098300706397494"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098388712332598"},
			BVIListIDs:            []string{"1829098388709186870"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098388711284022"},
			SpecialQuotaListIDs:   []string{"1829098388710235446"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098693184199990"},
			BVIListIDs:            []string{"1829098693180005686"},
			TargetQuotaListIDs:    []string{"1829184113205058870"},
			DedicatedQuotaListIDs: []string{"1829098693182102838"},
			SpecialQuotaListIDs:   []string{"1829098693181054262"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098438313123126"},
			BVIListIDs:            []string{"1829098438307880246"},
			TargetQuotaListIDs:    []string{"1829098438312074550", "1829185433435643190"},
			DedicatedQuotaListIDs: []string{"1829098438309977398"},
			SpecialQuotaListIDs:   []string{"1829098438308928822"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098789358542134"},
			BVIListIDs:            []string{"1829098789354347830"},
			TargetQuotaListIDs:    []string{"1829098789357493558"},
			DedicatedQuotaListIDs: []string{"1829098789356444982"},
			SpecialQuotaListIDs:   []string{"1829098789355396406"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098260805983542"},
			BVIListIDs:            []string{"1829098260798643510"},
			TargetQuotaListIDs:    []string{"1829098260801789238", "1829098260802837814", "1829098260803886390", "1829098260804934966"},
			DedicatedQuotaListIDs: []string{"1829098260800740662"},
			SpecialQuotaListIDs:   []string{"1829098260799692086"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098320139656502"},
			BVIListIDs:            []string{"1829098320134413622"},
			TargetQuotaListIDs:    []string{"1829098320138607926", "1829185325289708854"},
			DedicatedQuotaListIDs: []string{"1829098320136510774"},
			SpecialQuotaListIDs:   []string{"1829098320135462198"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098310219078966"},
			BVIListIDs:            []string{"1829098310214884662"},
			TargetQuotaListIDs:    []string{"1829184032826465590"},
			DedicatedQuotaListIDs: []string{"1829098310216981814"},
			SpecialQuotaListIDs:   []string{"1829098310215933238"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099000878341430"},
			BVIListIDs:            []string{"1829099000874147126"},
			TargetQuotaListIDs:    []string{"1829099000877292854"},
			DedicatedQuotaListIDs: []string{"1829099000876244278"},
			SpecialQuotaListIDs:   []string{"1829099000875195702"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098744493120822"},
			BVIListIDs:            []string{"1829098744484732214"},
			TargetQuotaListIDs:    []string{"1829098744487877942", "1829098744488926518", "1829098744489975094", "1829098744491023670", "1829098744492072246"},
			DedicatedQuotaListIDs: []string{"1829098744486829366"},
			SpecialQuotaListIDs:   []string{"1829098744485780790"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098799676529974"},
			BVIListIDs:            []string{"1829098799662898486"},
			TargetQuotaListIDs:    []string{"1829098799667092790", "1829098799668141366", "1829098799669189942", "1829098799670238518", "1829098799671287094", "1829098799672335670", "1829098799673384246", "1829098799674432822", "1829098799675481398", "1829185826818366774"},
			DedicatedQuotaListIDs: []string{"1829098799664995638"},
			SpecialQuotaListIDs:   []string{"1829098799663947062"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099071726427446"},
			BVIListIDs:            []string{"1829099071723281718"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099071725378870"},
			SpecialQuotaListIDs:   []string{"1829099071724330294"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098634136788278"},
			BVIListIDs:            []string{"1829098634128399670"},
			TargetQuotaListIDs:    []string{"1829098634132593974", "1829098634133642550", "1829098634134691126", "1829098634135739702", "1829185547806973238"},
			DedicatedQuotaListIDs: []string{"1829098634130496822"},
			SpecialQuotaListIDs:   []string{"1829098634129448246"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099090126839094"},
			BVIListIDs:            []string{"1829099090122644790"},
			TargetQuotaListIDs:    []string{"1829099090125790518"},
			DedicatedQuotaListIDs: []string{"1829099090124741942"},
			SpecialQuotaListIDs:   []string{"1829099090123693366"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098577307114806"},
			BVIListIDs:            []string{"1829098577302920502"},
			TargetQuotaListIDs:    []string{"1829185498868882742"},
			DedicatedQuotaListIDs: []string{"1829098577305017654"},
			SpecialQuotaListIDs:   []string{"1829098577303969078"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098399056534838"},
			BVIListIDs:            []string{"1829098399048146230"},
			TargetQuotaListIDs:    []string{"1829098399051291958", "1829098399052340534", "1829098399053389110", "1829098399054437686", "1829098399055486262"},
			DedicatedQuotaListIDs: []string{"1829098399050243382"},
			SpecialQuotaListIDs:   []string{"1829098399049194806"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098703560908086"},
			BVIListIDs:            []string{"1829098703548325174"},
			TargetQuotaListIDs:    []string{"1829098703551470902", "1829098703552519478", "1829098703553568054", "1829098703554616630", "1829098703555665206", "1829098703556713782", "1829098703557762358", "1829098703558810934", "1829098703559859510"},
			DedicatedQuotaListIDs: []string{"1829098703550422326"},
			SpecialQuotaListIDs:   []string{"1829098703549373750"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829116789938724150"},
			BVIListIDs:            []string{"1829116789925092662"},
			TargetQuotaListIDs:    []string{"1829116789928238390", "1829116789929286966", "1829116789930335542", "1829116789931384118", "1829116789932432694", "1829116789933481270", "1829116789934529846", "1829116789935578422", "1829116789936626998", "1829116789937675574"},
			DedicatedQuotaListIDs: []string{"1829116789927189814"},
			SpecialQuotaListIDs:   []string{"1829116789926141238"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098673703755062"},
			BVIListIDs:            []string{"1829098673698512182"},
			TargetQuotaListIDs:    []string{"1829098673701657910", "1829098673702706486"},
			DedicatedQuotaListIDs: []string{"1829098673700609334"},
			SpecialQuotaListIDs:   []string{"1829098673699560758"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098654396325174"},
			BVIListIDs:            []string{"1829098654393179446"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098654395276598"},
			SpecialQuotaListIDs:   []string{"1829098654394228022"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098663636376886"},
			BVIListIDs:            []string{"1829098663632182582"},
			TargetQuotaListIDs:    []string{"1829098663635328310"},
			DedicatedQuotaListIDs: []string{"1829098663634279734"},
			SpecialQuotaListIDs:   []string{"1829098663633231158"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098587239226678"},
			BVIListIDs:            []string{"1829098587235032374"},
			TargetQuotaListIDs:    []string{"1829185513459817782"},
			DedicatedQuotaListIDs: []string{"1829098587237129526"},
			SpecialQuotaListIDs:   []string{"1829098587236080950"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098448354286902"},
			BVIListIDs:            []string{"1829098448350092598"},
			TargetQuotaListIDs:    []string{"1829185446571154742"},
			DedicatedQuotaListIDs: []string{"1829098448352189750"},
			SpecialQuotaListIDs:   []string{"1829098448351141174"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098599124835638"},
			BVIListIDs:            []string{"1829098599113301302"},
			TargetQuotaListIDs:    []string{"1829098599116447030", "1829098599117495606", "1829098599118544182", "1829098599119592758", "1829098599120641334", "1829098599121689910", "1829098599122738486", "1829098599123787062"},
			DedicatedQuotaListIDs: []string{"1829098599115398454"},
			SpecialQuotaListIDs:   []string{"1829098599114349878"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098841639492918"},
			BVIListIDs:            []string{"1829098841635298614"},
			TargetQuotaListIDs:    []string{"1829098841638444342"},
			DedicatedQuotaListIDs: []string{"1829098841637395766"},
			SpecialQuotaListIDs:   []string{"1829098841636347190"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098611608132918"},
			BVIListIDs:            []string{"1829098611596598582"},
			TargetQuotaListIDs:    []string{"1829098611599744310", "1829098611600792886", "1829098611601841462", "1829098611602890038", "1829098611603938614", "1829098611604987190", "1829098611606035766", "1829098611607084342"},
			DedicatedQuotaListIDs: []string{"1829098611598695734"},
			SpecialQuotaListIDs:   []string{"1829098611597647158"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098270878604598"},
			BVIListIDs:            []string{"1829098270874410294"},
			TargetQuotaListIDs:    []string{"1829098270877556022"},
			DedicatedQuotaListIDs: []string{"1829098270876507446"},
			SpecialQuotaListIDs:   []string{"1829098270875458870"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098949899722038"},
			BVIListIDs:            []string{"1829098949894479158"},
			TargetQuotaListIDs:    []string{"1829098949897624886", "1829098949898673462"},
			DedicatedQuotaListIDs: []string{"1829098949896576310"},
			SpecialQuotaListIDs:   []string{"1829098949895527734"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099145912130870"},
			BVIListIDs:            []string{"1829099145907936566"},
			TargetQuotaListIDs:    []string{"1829099145911082294"},
			DedicatedQuotaListIDs: []string{"1829099145910033718"},
			SpecialQuotaListIDs:   []string{"1829099145908985142"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098459304566070"},
			BVIListIDs:            []string{"1829098459300371766"},
			TargetQuotaListIDs:    []string{"1829185459018800438"},
			DedicatedQuotaListIDs: []string{"1829098459302468918"},
			SpecialQuotaListIDs:   []string{"1829098459301420342"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098407870864694"},
			BVIListIDs:            []string{"1829098407867718966"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098407869816118"},
			SpecialQuotaListIDs:   []string{"1829098407868767542"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098683758550326"},
			BVIListIDs:            []string{"1829098683750161718"},
			TargetQuotaListIDs:    []string{"1829098683754356022", "1829098683755404598", "1829098683756453174", "1829098683757501750", "1829185560074263862"},
			DedicatedQuotaListIDs: []string{"1829098683752258870"},
			SpecialQuotaListIDs:   []string{"1829098683751210294"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098469564882230"},
			BVIListIDs:            []string{"1829098469557542198"},
			TargetQuotaListIDs:    []string{"1829098469560687926", "1829098469561736502", "1829098469562785078", "1829098469563833654"},
			DedicatedQuotaListIDs: []string{"1829098469559639350"},
			SpecialQuotaListIDs:   []string{"1829098469558590774"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829116820654660918"},
			BVIListIDs:            []string{"1829116820649418038"},
			TargetQuotaListIDs:    []string{"1829116820653612342", "1829185474661457206"},
			DedicatedQuotaListIDs: []string{"1829116820651515190"},
			SpecialQuotaListIDs:   []string{"1829116820650466614"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098490106486070"},
			BVIListIDs:            []string{"1829098490094951734"},
			TargetQuotaListIDs:    []string{"1829098490098097462", "1829098490099146038", "1829098490100194614", "1829098490101243190", "1829098490102291766", "1829098490103340342", "1829098490104388918", "1829098490105437494"},
			DedicatedQuotaListIDs: []string{"1829098490097048886"},
			SpecialQuotaListIDs:   []string{"1829098490096000310"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099010341739830"},
			BVIListIDs:            []string{"1829099010336496950"},
			TargetQuotaListIDs:    []string{"1829099010340691254", "1829185998784830774"},
			DedicatedQuotaListIDs: []string{"1829099010338594102"},
			SpecialQuotaListIDs:   []string{"1829099010337545526"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099029347179830"},
			BVIListIDs:            []string{"1829099029342985526"},
			TargetQuotaListIDs:    []string{"1829184190090845494"},
			DedicatedQuotaListIDs: []string{"1829099029345082678"},
			SpecialQuotaListIDs:   []string{"1829099029344034102"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099056144588086"},
			BVIListIDs:            []string{"1829099056141442358"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099056143539510"},
			SpecialQuotaListIDs:   []string{"1829099056142490934"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098523302305078"},
			BVIListIDs:            []string{"1829098523298110774"},
			TargetQuotaListIDs:    []string{"1829098523301256502"},
			DedicatedQuotaListIDs: []string{"1829098523300207926"},
			SpecialQuotaListIDs:   []string{"1829098523299159350"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098972660112694"},
			BVIListIDs:            []string{"1829098972643335478"},
			TargetQuotaListIDs:    []string{"1829098972646481206", "1829098972647529782", "1829098972648578358", "1829098972649626934", "1829098972650675510", "1829098972651724086", "1829098972652772662", "1829098972653821238", "1829098972654869814", "1829098972655918390", "1829098972656966966", "1829098972658015542", "1829098972659064118"},
			DedicatedQuotaListIDs: []string{"1829098972645432630"},
			SpecialQuotaListIDs:   []string{"1829098972644384054"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099063853718838"},
			BVIListIDs:            []string{"1829099063851621686"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099063852670262"},
			SpecialQuotaListIDs:   []string{},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099019660434742"},
			BVIListIDs:            []string{"1829099019656240438"},
			TargetQuotaListIDs:    []string{"1829184175621545270"},
			DedicatedQuotaListIDs: []string{"1829099019658337590"},
			SpecialQuotaListIDs:   []string{"1829099019657289014"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829116694573882678"},
			BVIListIDs:            []string{"1829116694572834102"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{},
			SpecialQuotaListIDs:   []string{},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099038427847990"},
			BVIListIDs:            []string{"1829099038424702262"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829099038426799414"},
			SpecialQuotaListIDs:   []string{"1829099038425750838"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098500119338294"},
			BVIListIDs:            []string{"1829098500110949686"},
			TargetQuotaListIDs:    []string{"1829098500115143990", "1829098500116192566", "1829098500117241142", "1829098500118289718", "1829185487097568566"},
			DedicatedQuotaListIDs: []string{"1829098500113046838"},
			SpecialQuotaListIDs:   []string{"1829098500111998262"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098330035068214"},
			BVIListIDs:            []string{"1829098330029825334"},
			TargetQuotaListIDs:    []string{"1829098330034019638", "1829185341341310262"},
			DedicatedQuotaListIDs: []string{"1829098330031922486"},
			SpecialQuotaListIDs:   []string{"1829098330030873910"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098851832700214"},
			BVIListIDs:            []string{"1829098851824311606"},
			TargetQuotaListIDs:    []string{"1829098851828505910", "1829098851829554486", "1829098851830603062", "1829098851831651638", "1829185851031035190"},
			DedicatedQuotaListIDs: []string{"1829098851826408758"},
			SpecialQuotaListIDs:   []string{"1829098851825360182"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829116768661019958"},
			BVIListIDs:            []string{"1829116768652631350"},
			TargetQuotaListIDs:    []string{"1829116768656825654", "1829116768657874230", "1829116768658922806", "1829116768659971382", "1829185894060399926"},
			DedicatedQuotaListIDs: []string{"1829116768654728502"},
			SpecialQuotaListIDs:   []string{"1829116768653679926"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098871332019510"},
			BVIListIDs:            []string{"1829098871325728054"},
			TargetQuotaListIDs:    []string{"1829098871329922358", "1829098871330970934", "1829185863905451318"},
			DedicatedQuotaListIDs: []string{"1829098871327825206"},
			SpecialQuotaListIDs:   []string{"1829098871326776630"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098881359551798"},
			BVIListIDs:            []string{"1829098881352211766"},
			TargetQuotaListIDs:    []string{"1829098881356406070", "1829098881357454646", "1829098881358503222", "1829185933242539318"},
			DedicatedQuotaListIDs: []string{"1829098881354308918"},
			SpecialQuotaListIDs:   []string{"1829098881353260342"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098919855922486"},
			BVIListIDs:            []string{"1829098919852776758"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{"1829098919854873910"},
			SpecialQuotaListIDs:   []string{"1829098919853825334"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098891417492790"},
			BVIListIDs:            []string{"1829098891410152758"},
			TargetQuotaListIDs:    []string{"1829098891414347062", "1829098891415395638", "1829098891416444214", "1829185876179033398"},
			DedicatedQuotaListIDs: []string{"1829098891412249910"},
			SpecialQuotaListIDs:   []string{"1829098891411201334"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098339801505078"},
			BVIListIDs:            []string{"1829098339797310774"},
			TargetQuotaListIDs:    []string{"1829185370765401398"},
			DedicatedQuotaListIDs: []string{"1829098339799407926"},
			SpecialQuotaListIDs:   []string{"1829098339798359350"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098900757159222"},
			BVIListIDs:            []string{"1829098900751916342"},
			TargetQuotaListIDs:    []string{"1829098900756110646", "1829185946062429494"},
			DedicatedQuotaListIDs: []string{"1829098900754013494"},
			SpecialQuotaListIDs:   []string{"1829098900752964918"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098535633558838"},
			BVIListIDs:            []string{"1829098535627267382"},
			TargetQuotaListIDs:    []string{"1829098535630413110", "1829098535631461686", "1829098535632510262"},
			DedicatedQuotaListIDs: []string{"1829098535629364534"},
			SpecialQuotaListIDs:   []string{"1829098535628315958"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098910436564278"},
			BVIListIDs:            []string{"1829098910431321398"},
			TargetQuotaListIDs:    []string{"1829098910435515702", "1829185960242322742"},
			DedicatedQuotaListIDs: []string{"1829098910433418550"},
			SpecialQuotaListIDs:   []string{"1829098910432369974"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098734183521590"},
			BVIListIDs:            []string{"1829098734171987254"},
			TargetQuotaListIDs:    []string{"1829098734175132982", "1829098734176181558", "1829098734177230134", "1829098734178278710", "1829098734179327286", "1829098734180375862", "1829098734181424438", "1829098734182473014"},
			DedicatedQuotaListIDs: []string{"1829098734174084406"},
			SpecialQuotaListIDs:   []string{"1829098734173035830"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098821050703158"},
			BVIListIDs:            []string{"1829098821041265974"},
			TargetQuotaListIDs:    []string{"1829098821045460278", "1829098821046508854", "1829098821047557430", "1829098821048606006", "1829098821049654582", "1829185838317051190"},
			DedicatedQuotaListIDs: []string{"1829098821043363126"},
			SpecialQuotaListIDs:   []string{"1829098821042314550"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098811032608054"},
			BVIListIDs:            []string{"1829098811022122294"},
			TargetQuotaListIDs:    []string{"1829098811025268022", "1829098811026316598", "1829098811027365174", "1829098811028413750", "1829098811029462326", "1829098811030510902", "1829098811031559478"},
			DedicatedQuotaListIDs: []string{"1829098811024219446"},
			SpecialQuotaListIDs:   []string{"1829098811023170870"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829116745423527222"},
			BVIListIDs:            []string{"1829116745413041462"},
			TargetQuotaListIDs:    []string{"1829116745416187190", "1829116745417235766", "1829116745418284342", "1829116745419332918", "1829116745420381494", "1829116745421430070", "1829116745422478646"},
			DedicatedQuotaListIDs: []string{"1829116745415138614"},
			SpecialQuotaListIDs:   []string{"1829116745414090038"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829098418061974838"},
			BVIListIDs:            []string{"1829098418052537654"},
			TargetQuotaListIDs:    []string{"1829098418055683382", "1829098418056731958", "1829098418057780534", "1829098418058829110", "1829098418059877686", "1829098418060926262"},
			DedicatedQuotaListIDs: []string{"1829098418054634806"},
			SpecialQuotaListIDs:   []string{"1829098418053586230"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829099047875517750"},
			BVIListIDs:            []string{"1829099047871323446"},
			TargetQuotaListIDs:    []string{"1829099047874469174"},
			DedicatedQuotaListIDs: []string{"1829099047873420598"},
			SpecialQuotaListIDs:   []string{"1829099047872372022"},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829101006431984950"},
			BVIListIDs:            []string{"1829101006430936374"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{},
			SpecialQuotaListIDs:   []string{},
		},

		&mirea.HTTPHeadingSource{
			RegularListIDs:        []string{"1829100731190222134"},
			BVIListIDs:            []string{"1829100731189173558"},
			TargetQuotaListIDs:    []string{},
			DedicatedQuotaListIDs: []string{},
			SpecialQuotaListIDs:   []string{},
		},
	}
}
