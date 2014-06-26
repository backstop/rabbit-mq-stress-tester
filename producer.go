package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type ProducerConfig struct {
	Uri        string
	Bytes      int
	Quiet      bool
	WaitForAck bool
}

func Produce(config ProducerConfig, tasks chan int) {
	connection, err := amqp.Dial(config.Uri)
	if err != nil {
		println(err.Error())
		panic(err.Error())
	}

	channel, err1 := connection.Channel()
	if err1 != nil {
		println(err1.Error())
		panic(err1.Error())
	}

	if config.WaitForAck {
		channel.Confirm(false)
	}

	ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	q := MakeQueue(channel)

	for {

		sequenceNumber, alive := <-tasks

		if !alive {
			channel.Close()
			connection.Close()
			log.Println("Broke out of loop!")
			return
		}

		start := time.Now()

		message := &MqMessage{start, sequenceNumber, makeString(config.Bytes)}
		messageJson, _ := json.Marshal(message)

		channel.Publish("", q.Name, true, false, amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "UTF-8",
			Body:            messageJson,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
		)

		confirmOne(ack, nack, config.Quiet, config.WaitForAck)

		if !config.Quiet {
			log.Println(time.Since(start))
		}
	}

}

func confirmOne(ack, nack chan uint64, quiet bool, waitForAck bool) {
	if waitForAck {
		select {
		case tag := <-ack:
			if !quiet {
				log.Println("Acked %d", tag)
			}
		case tag := <-nack:
			log.Println("Nack alert! %d", tag)
		}
	}
}

func makeString(bytes int) string {
	longString := `
	ANTAŬPAROLO.


Estas afero konata, ke la malgrandaj nacioj tre ŝatas esti
fieraj pro siaj grandaj homoj kaj ke ili tre volonte atentigas la
internacian publikon pri la fakto, ke la penso ekbruliginta korojn kaj
kapojn de homaro devenis de filo de ilia tero kaj ilia nacio.

La instruo pri la internacia solidareco, estu ĝi disvastiĝinta
kiom ajn, neniam forpelos el homaj koroj la sanktan flamon de speciala
amo al homoj naskiĝintaj en la sama regiono kaj lernintaj de la
patrino la saman idiomon, amon de la propra nacio. Ĉar la kaŭzo de tiu
ĉi amo estas la plej esenca kvalito de homaro: ĝi estas kredo je forto
de la homa penso, je povo per la homecaj ideoj regi la sorton de si
mem kaj sorton de la tuta homaro.

Kaj la plej altaj, la plej idealaj kaj sekve la plej efikaj
ideoj progresigi la homaron estas tiuj de mi, simpla, senpova kaj
sensignifa membro de sia nacio, sed forta tiam, se koroj de ĉiuj miaj
samnacianoj ekbatas por la samaj savontaj ideoj.

Kaj kun ĝojega miro mi aŭdas, ke, kio estis idealo kaj savonta
vorto por mi, estas same idealo kaj devizo de la ceteraj samnacianoj;
eĉ ke la ideoj de la ceteraj estas pli klaraj kaj ke la ilia koro
batas pro ili ankoraŭ pli brue ol la mia; kaj ke ekzistas inter la
samnacianoj homoj, kiuj scias diri miajn idealojn per tiaj vortoj kaj
montri al ili tiel ireblajn vojojn, ke la kapo flarante poveblecon
realigi idealojn pensitajn nerealigeblaj svenas pro feliĉego. Tio
estas la tasko de la grandaj animoj. Ĉu do oni ne devas fieriĝi pro
nacia parenceco kun la granda spirita korifeo? Ĉu estas honto tia
fiero?

Jen la argumento, kial la ĉeĥoslovaka nacio fieriĝas pro Jan
Amos Komenský (Comenius) kaj kial ni eldonas en esperanta traduko
broŝuron priskribantan liajn vivon kaj laboron.

Jan Amos Komenský estas speciale kara al koro de ĉiu
Ĉeĥoslovako, ĉar li estas unu de la tri homoj, kiuj plej esence
reprezentas internan ideon de la ĉeĥoslovaka nacio, sencon de ĝia
historio, kiun dum la lastaj jardekoj tre kuraĝe esploradis niaj
filozofoj kaj historiistoj precipe T. G. Masaryk (nuna prezidanto de
la respubliko). Lia teorio trovis multan aprobon kaj multenombrajn
kredantojn. Jen ĝia esenco.

Senco kaj celo de ekzisto de la ĉeĥoslovaka nacio, la interna
ideo de ĝia evoluo, estas montrita en ĝia historio. El ĝi sekvas, ke
la nacio estis en la plej fruktodona evoluo de siaj kapablecoj tiam,
kiam ĝi luktis ekskluzive por celoj abstraktaj, spiritaj, idealaj,
kiuj ne povis efiki je lia materia situacio. Tio estis dum epoko de
Hus (bruligita en j. 1415) kaj de Husanoj.

La parolo de Hus „Serĉu la veron kaj restu en la konita vero!“
kiu estis de Masaryk montrita kiel la precipa instigilo de la Husa
celado kaj fariĝis devizo de la progresema movado, estas mallongigo de
la plej arde dezirata idealo nacia: povi klopodi por koni la veron kaj
por efektivigi postulojn kaj ordonojn de la vero en la praktika
vivo. Kio ĉe tio estas speciale substrekinda estas la efektivigo de la
vero en la praktiko, penado, ke la ĉiutaga, nefestena, laborplena vivo
fariĝu enkorpiĝo, realiĝo, efektiviĝo de la verecaj idealoj. La
konkretaj faroj povas esti eĉ kontraŭaj, sed tiu ĉi pensmaniero estas
komuna al ĉiuj grandaj aperaĵoj de la spirita vivo ĉeĥoslovaka.

La dua fiero de nia historio estas Petr Chelčický, samtempulo
de Hus kaj de la Husanoj. La Husanoj decidis batali por sia vero, por
sia vero de Dio, ili entuziasmigis por la „Dia milito“ la plimulton de
la nacio, ili jam en jaro 1420 ricevis gravan venkon kontraŭ la neleĝa
reĝo Zigmundo. Sed tiam Chelčický havas kuraĝon kontraŭstari la
ĝeneralan humoron de la tempo kaj ekpredikas opinion elpensitan en
silento de la sudbohemaj kampoj: „Ĉia ajn milito estas malpermesita al
Kristano per la unusenca ordono Dia: Nemortigu! Neniam decas al
Kristano uzi la perforton“. Kaj kiam la ĉefurbo plena de esperoj je
militaj venkoj kaj venkaj profitoj ne havis orelon por aŭskulti lin,
li revenis en sian kamparon, ne por verki simple novajn librojn--sed
por organizi homojn, kiuj ekkredinte liajn verojn estis kapablaj vivi
laŭ ili. Kaj li efektivigis tion kun tiela sukceso, ke fondiĝis multaj
komunumoj de homoj vivantaj laŭ la principoj de la origina kristaneco,
kiel ĝin Chelčický predikis. Ili neniam uzis perforton, amis sin
reciproke, laboris honeste ne sole por sia propra profito sed por
bonfarto de ĉiuj komunumanoj. Ilia nomo estis „Ĉeĥaj Fratoj“ kaj la
influo de organiza laboro de „revulo“ Chelčický daŭris tri jarcentojn!
Dum tiel longa epoko kreskis kaj floris komunumoj de Ĉeĥaj Fratoj kiel
modeloj de la bona mastrumado kaj dokumento, ke eĉ utopioj estas
realigeblaj.

Jan Amos Komenský estis la lasta episkopo de la Ĉeĥaj
Fratoj. Sed lia vero estis jam alia vero. La mezepokaj formoj estas ĉe
li jam tute plenigitaj per la moderna enhavo de renesanca kulturo, la
du polusoj de la homara celado: idealo materia kaj idealo spirita, la
bezono agi kaj la ĝuo mediti, estas ĉe li en harmonia kunsento. Li
estas efektive kredanta membro de la eklezio de la Ĉeĥaj Fratoj, li
malbenas kaj kondamnas militon, sed antaŭ ĉio li scias, ke lia propra
nacio devas esti libera por povi plenumi taskojn de Dio
donatajn. Samtempe kun li vivis la glora filozofo franca
Descartes. Lia pens- kaj agadmaniero estis tute kontraŭa al tiu de
Komenský.

Descartes dum la tuta vivo sidis en salono sur kusenoj
malforta, malsaneta viro, kaj la tutan vivon meditis, ĉu li ekzistas
aŭ ne.

Al Komenský eĉ ne venis penseto, ke li povus ne ekzisti, ĉar
ĉiam aferoj gravegaj estis farotaj. Li, politika ekzilito el sia
patrujo, ĉiam sciis, ke antaŭ ĉio la maljusta kaj malbona uzurpatoro
de la reĝa trono estas forigota. Kaj li tial vagadis tra la tuta
Eŭropo, serĉante monon, soldatojn, asociojn, kiuj povebligus forigi la
malhonestulon Habsburgan, la despoton, la kontraŭreformaciulon. Eĉ kun
Turkoj li estus estinta asociinta por ekmiliti la Habsburgojn.

Kaj kiam ĉio montriĝis vana, li ne malesperis sed ĉiujn siajn
laborkapablojn li dediĉis al la reformo de pedagogio al la organizado
de pli bonaj lernejoj ol nun. Ĉar se la junaj homoj eliros el lernejoj
pli bone instruitaj pri la homecaj virtoj, tiam neniam estos poveble,
ke ekzistu tiomaj maljustecoj, tiom da faroj kontraŭ la leĝo de Dio,
kiom li vidis kaj travivis.

Tia konkreta, praktika laboro por la ideala vero, por la
morala postulo, por sistemo de Dio, tio estas la plej karakteriza
kvalito de la ĉeĥoslovaka menso.

Tial komprenu, fremdnaciaj legantoj de nia libreto, kial ni
kun tia fervoro propagandas la scion pri Jan Amos Komenský, kaj helpu
nin propagandante kun ni.

BRATISLAVA, 13. marton 1921.

Dr. Stan. Kamaryt.



JAN AMOS KOMENSKÝ,
INSTRUISTO DE LA NACIOJ, FONDINTO
DE LA NOVOEPOKA INSTRUADO.


„Ankaŭ mi kredas al Dio, ke ĝis pasos la uragano de la kolero
pro niaj pekoj sendita sur niajn kapojn, la regado de viaj aferoj al
vi revenos, ho, popolo bohema!“

Per tiuj ĉi vortoj, ofte nun ripetataj, inaŭguraciis la unua
prezidanto de la ĉeĥoslovaka respubliko Tomáš G. Masaryk sian unuan
komunikon al la Nacia Kunveno. Kaj la profeto Jan Amos Komenský, kies
plumo ilin la unuan fojon skribis, tiam ilin skribis, kiam montriĝis
vantaj la lastaj esperoj de bohemaj ekzilitoj je politika helpo, post
ratifikacio de la interkonsentoj de paco Westfala en Nürenberg
(komencon de la jaro 1650) kaj kiam la turmentata eŭropa popolo denove
libere ekspiris post la longa suferado pro la plej kruela eklezia
milito. Por li kaj liaj samreligianoj tiu trankvilo bedaŭrinde ne
ekzistis, al ili estis destinite daŭrigi la vagadon tra la negastema
fremdlando, kaj trovi tie eĉ la finan ripozon.

Kaj tamen Komenský, la lasta ĉefpastro de la "Frata Unuiĝo"
(Jednota bratrská), nia plej nacia eklezio, senlace laboris preskaŭ
duonon da jarcento por la bono de la homaro dediĉante la amon al sia
patrujo, kiam li devis ĝin forlasi en aĝo da 35 jaroj, penante ĝis sia
lasta spiro helpi al ĝi el la morala mizero per disvastiĝo de la
klereco en ĉiujn klasojn kaj per la konvinkado pri bezono de la paco.

1. Riĉa je sciaĵoj li revenis el fremdlando, por dediĉi ĉiujn
siajn fortojn al la servo de la Frata Eklezio; sed jam kiel studento
li ekkuraĝis konatigi al siaj samlandanoj elpensaĵojn de la samtempaj
natursciencoj kaj per helpo de la detala ĉeĥolatina vortaro ebligi al
ili la scion de la sciencoj.

Sed bedaŭrinde lia trankvila laboro por la utilo de liaj
eklezianoj estis interrompita per eksplodo de ribelo, finiĝinta per
venko de la reĝo Ferdinando IIa super la bohema nacio la 8an novembron
1620.

En tiu tempo estis ankaŭ verkita lia fama „Labirinto de la
mondo kaj paradizo de la koro“ (1623), kiun libron oni ankoraŭ hodiaŭ
vicigas inter la popularajn librojn de la pli malnova ĉeĥa literaturo.

2. Sed ankaŭ pri la estonteco de la patrujo zorgis Komenský
dum ĉi tiuj suferoplenaj jaroj. Li estis plene konvinkita, ke ne estas
malproksima la tempo, kiam finos la suferado pro la kredo kaj ekestos
la nova regno, pri kiu ĉiutage preĝas Kristanoj per vortoj de la preĝo
„Patro nia. Venu regno Via“! Zorgante tiel, ke lia nacio ne estu
surprizita per la subita renverso, li preparis por ĝi proponon de la
nova reguligo de la eklezio kaj samtempe la verkon, kiu garantiis lian
gloron por la tuta estonteco. La verkon li nomis „Didaktiko,“
t. e. „arto de la arta instruado“ kaj en ĝi li montris pli frue al la
samlandanoj, pli poste en la latina traduko de la verko al la tuta
homaro, kiamaniere eduki la junularon al la noblaj moroj kaj pieco kaj
instrui al ĝi utilajn sciojn. En la libro oni postulas dediĉi por la
perfekta eduko plenajn 24 jarojn de la juneca aĝo, kiun oni dividu je
kvar periodoj; la infanaĝo, knabaĝo, aĝo de la maturiĝado kaj
junulaĝo. Por ĉiu aĝo li volis starigi apartan lernejon, preparante
por ĉiu specialajn librojn. Al la infana aĝo li dediĉis la faman
„Informatorio de la patrina lernejo“, kiu fariĝis evangelio de la
instruistaro ne nur de lia propra nacio, sed ankaŭ de ĉiuj kleraj
nacioj same kiel la „Didaktiko“. Bedaŭrinde ĉi tiujn signifoplenajn
fruktojn de lia spirito povis utiligi pli frue la fremdlando, liaj
samlandanoj ekuzis ilin nur post la nacia renaskiĝo.

3. En ekzilon en la urbon Leszno foriris Komenský je komenco
de la jaro 1628 kaj tie li komencis instrui lernantojn de la latina
lernejo. Li verkis por ili negrandan lernolibron de la latina lingvo
nomata „Janua linguarum reserata“ (1631), kiu disvastiĝis dum kelkaj
jaroj en la tuta Eŭropo kaj konatigis la nomon de Komenský eĉ en la
plej malproksimaj landoj de Azio. Ĝia praktika aranĝo ekinteresis la
tutan instruitan mondon. Komenský deziris per tiu ĉi lernosistemo ne
sole forigi malfacilaĵojn de la lernado de la lingvo latina, sed li
celis samtempe doni al la junularo konon de la plej gravaj aferoj en
la mondo. Tial li konstruis en 100 alineoj (1000 frazoj) tabelon de
enhavo de la tuta mondo, komencante de ĝia kreo kaj finante per la
lastaj aferoj de la homo; en tio ĉiu objekto devis esti signita per
sia origina nomo. La praktika signifo de tia libro estis al ĉiu
evidenta kaj tial ĝi estis en la tuta Eŭropo en diversajn lingvojn
tradukata. La verkinto mem ekmiris pri la sukceso kaj aprobo de sia
libro. La provo faris lin tuj mondkonata inter ĉiuj tiamaj
kleruloj. Sed tamen en sia fervora penado helpi al la homaro,
skuiĝanta pro la malpaco kaj maltrankvilo, li konsideris tion kvazaŭ
signon de la Providenco, ke li provu ion pli signifan. Li rimarkis, ke
kaŭzo de ĉiuj malfacilaĵoj sur la mondo estas celado al aferoj
supraĵaj sen atento de ilia interno kaj enhavo kaj la homa
facilkredemeco. Antaŭ nelonge estis atendata savo de la klasikismo, la
reveno al beleco de la malnovaj formoj literaturaj, sed la penado
haltis ĉe senenhavaj formoj, ĉar al la efektiva kono, ne alkondukis la
junularon eĉ la legado de la malnovaj verkistoj. Kaj tamen estis en la
lernejoj, tiom da tempo kaj zorgoj dediĉita al la latina lingvo sole
por ke la lernantoj sciu verki leterojn kaj versojn laŭ la malnovaj
modeloj. Ĉu estas utila por la vivo, demandas Komenský, labori pri
tiaj aferoj, kiam ĵus dume atentigas pri sia rajto nove prosperantaj
natursciencoj precipe fiziko, ĥemio kaj astronomio? Tial opinias
Komenský, ke estas necese aranĝi denove la instruadon laŭ novaj
metodoj kaj principoj, kiel oni faris ĉe lingvo latina. Estas jam por
la komparo de ĉiuj objektoj necese tiujn principojn trovi. Ke estas
eble fari tion, konkludas Komenský laŭ vortoj de Sankta Skribo (Saĝeco
11, 21), ke Dio ĉion ordigis laŭ la nombro, mezuro kaj pezo. Se ĉiuj
vortoj de iu lingvo havas sian analogion, ke oni povas ilin ordigi en
nemultaj reguloj, sammaniere klasifiki oni povas la objektojn kaj
difini laŭ iom da principoj. Estus sole nur necese konstrui „Pordegon
de la objektoj“, memkompreneble tiel konstruitan, ke la junularo sen
malfacilaĵoj eniru en la templon de la sciencoj. Kaj se estus tiel
poveble de la juneco ĉiujn homojn enkonduki en la veran konon,
malaperus certe per si mem ĉiuj iniciatoj de la malkonsento kaj
malpaco kaj anstataŭ malamikeco ekregus en la mondo agrabla paco. Kaj
pri paco klopodis Komenský des pli multe kun aliaj same pensantaj
viroj, ju pli profunde influis la ĉeĥan landon, lian patrujon, kaj
preskaŭ la tutan Eŭropon la malfeliĉaj sekvoj de la religiaj
malpacoj. Se la vera scio disvastiĝus inter ĉiuj homoj, ĉu ne
proksimiĝus la plenumo de la vortoj de Kristo, ke estos unu ŝafejo kaj
unu paŝtisto?

5. Ke oni iam sukcesos verki tian libron, Komenský ne
dubis. Ankaŭ ĝustan nomon li el libroj samtempaj por ĝi trovis:
„Pansofio“. Kompreneble por ke ĉiuj frazoj en ĝin akceptitaj estu vere
fidindaj konitaĵoj, oni devus trastudi tutan sciencan materialon kaj
per spertoj kaj esplorado pruvi la verecon de ĉiu ĝia frazo. Sed tio
estus granda kaj malfacila laboro, ke sole senlaca kaj diligenta
laboro de unuopulo sukcesus ĝin venki.

6. La ideojn, en ĝi koncize donitajn, pritraktis Komenský
detale en la letero, kiun li sendis al sia amiko Samuel Hartlib en
London, en la duono de tridekaj jaroj de la XVII. jarcento. Al Hartlib
ĉio ĉi tre plaĉis, ĉar ankaŭ la fama angla filozofo lordo Francis
Bacon, antaŭ ne longe mortinta (1628), predikis la bezonon de nova
metodo de ĉia esplorado sur la principo de indukcio, pri kiu tiam
klopodis la natursciencoj. Estas tiel klarigeble, ke Hartlib sen
permeso de Komenský sendis la leteron al la Oxforda universitato por
bontrovo kaj kun ĝia konsento li ĝin eĉ prese eldonis kiel „Conatuum
Comenianorum praeludia“ (Antaŭludo de entreprenoj de
Komenský). Komenský estis vere surprizita ricevinte subite el Londono
anstataŭ respondo presitajn ekzemplerojn de sia
letero. Memkompreneble, ke ne tute konvenis al Komenský antaŭ la mondo
entrepreni ion, kion li ŝatus labori kaŝe en sia laborejo, sed
kuraĝigis lin denove la konsento kun lia propono de la tuta scienca
mondo. Baldaŭ estis proponite al li ne nur la materiala helpo por lia
laboro, sed ankaŭ la scienculoj de diversaj fakoj proponis sian helpon
kiel kunlaborantoj, por ke la verko estu pli frue preta.

7. Unuan materialan helpon proponis al Komenský la juna
posedanto de Leszno Grafo Bohuslav Leszczynski, kiu volis ankaŭ liajn
helpantojn vivteni. Tio okazis en la jaro 1640, post granda sukceso de
la lerneja teatraĵo „Pri Diogenes el Kynike“, kiun li atingis estante
rektoro de la leszna lernejo. La grafo ĉeestis kun siaj amikoj la
prezentadon. Al li do Komenský prezentis unuan detalan proponon, kiel
estus eble pansofion sistemi kaj kiuj helpaj verkoj estas bezonaj. Ĝi
estis precipe „Panhistorio“, historio de scienca esplorado kaj serĉado
de la veroj kaj la „Pandogmatio“, resumo de ĉiuj opinioj kaj skismoj
per kiuj la homaro iradis al la scio.

8. Sed pli frue ol ĉi tiu propono estis decidita, ricevis
Komenský de Hartlib (somero 1641) inviton, ke li venu Anglion por
signifoplena laboro, gloranta kaj honoranta Dion. Alveninte Londonon
Komenský nur tiam eksciis de la amikoj, ke li estis invitita de la
parlamento mem, por fondi pansofian kolegion, en kiu li kun 12
scienculoj laborus je la progreso de la sciencoj. Tiun sukceson
atingis Hartlib helpe de potencaj amikoj ĉe la nove kunvokitaj
parlamentanoj, ĉe kiuj tiam aperis forte la emo, ne nur plibonigi
ekzistantajn malordojn en la regno, sed ankaŭ subteni la progresemajn
ideojn, kiuj tiam estis sennombraj. Bedaŭrinde, la ribelo de katolikaj
Irlandanoj kontraŭ Angloj prokrastis la tujan efektivigon de la
grandioza propono de Hartlib.

9. Komenský uzis liberan tempon en Londono dum vintro 1641-2
preparante detalan planon nomitan „Via Lucis“ (Vojo de la lumo) por la
intencita kolegio. Li montras en ĝi kiamaniere penis la homaro
liberiĝi el la mallumo de la erardoktrinoj, en kiuj ĝi kvazaŭ droniĝis
post perdo de la kredo je la vera Dio. Lastaj kaj plej efikaj rimedoj,
laŭ opinio de Komenský, estas la arto de libropresado kaj vojaĝado
trans maroj. Sed ankaŭ ĉi tio ne sufiĉas por ke la homaro retrovu la
lumon. Por atingi ĝin oni ĉefe bezonas disvastigi la klerecon per
taŭgaj rimedoj. Tiuj rimedoj povas esti elementaj lernejoj por la tuta
homaro, en kiuj ĝi akiradus la konon el universalaj libroj, verkitaj
de scienca kolegio. Ankaŭ universalan lingvon povus iam tiuj ĉi viroj
krei. Havante sian centron en Londono, kvazaŭ foirejo de la tuta mondo
komerca, ili povus atente observi la evoluon de la tutmonda kulturo,
informi sin reciproke pri novaj eltrovoj de monda signifo kaj esti
samtempe internaciaj juĝantoj en aferoj politikaj. Al ilia juĝorajto
estus ankaŭ la reĝoj submetataj kaj kiu kontraŭstarus ilin, tiu povus
esti forigita per la komuna interveno de ĉiuj aliaj. Tiamaniere estus
garantiita la senĉesa progresado de la homaro kaj la eterna paco en la
tuta mondo.

10. Bedaŭrinde, la grandioza projekto ne efektiviĝis, la
politika paco ne ekestis en Anglio eĉ dum la sekvintaj jaroj kaj la
ribelo kontraŭ la malnovaj ordoj ekkulminaciis per forigo de la reĝeco
kaj enkonduko de protektorato de Cromwell. Ankaŭ religia radikala
partio sukcesis valori en la politiko. Komenský povis tiam nenion
alian fari ol forlasi konfuzitan Anglion kaj akcepti novan proponon el
Svedio, ke li organizu por Svedoj la librojn por lernado de latina
lingvo kaj kompense ke li ricevos helpon por la preparata pansofio.

11. El simpla „Janua linguarum reserata“ evoluis dum kelke da
jaroj triparta metodo de latinaj lernolibroj, kiuj estis sukcese
enkondukataj en tiamajn lernejojn. Krom la „Janua linguarum“ estis
ankoraŭ „Antaŭpordejo“ (Vestibulum) kaj „Ĉambro“ (Atrium) al ĉiu grado
de la gramatiko kaj vortaro. Komenský nevolonte forlasis ŝatatajn jam
ideojn pansofiajn, en kiuj li vidis sian propran vivocelon, precipe,
kiam liaj anglaj amikoj lin pro tio riproĉis, sed li esperis baldaŭ
fini la laboron en Svedio kaj komenci tiun, kiu plej multe lin
allogis. Ankaŭ pro politikaj kaŭzoj li volis plenumi deziron de
Svedoj, kiuj kiel la plej potencaj kontraŭloj de Habsburgoj penis
helpi liajn samlandanojn kaj samkredanojn kaj ebligi la revenon en
Bohemion.

12. Trans Holandio, kie li ankaŭ renkontiĝis kun la fama
filozofo Descartes, Germanio kaj Danio li rapidis Svedion, por
prezenti sin al la malavara subtenanto de subpremataj homoj, al
holanda grandkomercisto Ludoviko de Geer, militistara sveda liveranto,
de kiu la propono devenis. Tiu, volante havi sian intencon aprobitan
de la politika potenco, sendis lin en Stockholmon al regna
grandkanceliero Axel Oxenstern. La kanceliero, jam de longe konscianta
bezonon de la reformo de la instruado, aŭskultis la konkludojn de
Komenský kun granda intereso, li eĉ proponis al li superan inspekton
de reformo de la tuta instruado en Svedio. Ne malpli ĝojigadis lin,
kiam Komenský argumentis al li el la Sankta Skribo la proksimecon de
la regno de Kristo kun la pligrandigo de komuna kulturo. Nur post la
kvartaga interparolado foriris Komenský, promesinte, ke li pli frue,
kiel aferon pli urĝan, prilaboros librojn lingvajn, kaj poste komencos
solvadi la pansofion (Ĉioscion). Ankaŭ al reĝino Kristino estis tiam
Komenský de kanceliero prezentita.

Kiel loĝlokon por la komencotaj laboroj, rekomendis al li la
kanceliero trankvilan, tolereman urbon Elbing, apud freŝa limano de la
rivero Vistulo, en la orienta Prusio. Tie ankaŭ ekloĝis Komenský kun
sia familio, en aŭtuno de la jaro 1642, havante kelkajn junajn
helpantojn el Frata Unuiĝo, kun kiuj li esperis sukcese plenumi la
tutan taskon dum kelke da monatoj.

13. Sed denove prilabori tutan sistemon de lernolibroj kaj
samtempe meti ilin en perfektan harmonion, postulis multon da laboro
kaj pacienco precipe, kiam Komenský, dezirante perfekte koni la tutan
problemon, ree trastudis ĉiujn samtempajn lingvajn instrumetodojn kaj
eĉ analitike relaboris tutan instruan teorion entute, kaj nur sur ĉi
tiu fundamento li traktis, laŭ lia opinio, multsignifan por la tuta
klero teorion de ĝusta instruado de lingvoj en enkonduka verko al la
lernolibroj nomita „Linguarum methodus novissima“ (Plej nova lingva
metodo, eld. en jaro 1649).

14. Sed malgraŭ tiuj ĉi tiom urĝaj laboroj, li ne forgesadis
sian pansofion, eldoninte tuj en la jaro 1643 „Planon de pansofio“
(Pansophiae diatyposis) kaj 2 jarojn pli poste li komencis esplori la
penson de komuna universala reformo de ĉiu homa gento per sistemo de
la libroj, inter kiuj ankaŭ lia pansofio estus nur unusola verko. Li
intencis plej frue prezenti la proponon al Eŭropa klerularo por
prikonsidero kaj li nomis ĝin „Ĝenerala interkonsilo pri la reformo de
la aferoj homaj“ (De rerum humanarum emendatione consultatio
catholica). La efektiviĝo de tiu ĉi intenco ne okazis, ĉar Ludoviko de
Geer senĉese insistis je la finlaborado de la prenita tasko kaj li ne
kaŝis eĉ sian malplaĉon, se Komenský deflankiĝis for la laboro al tiu,
kiu lin pli allogis, aŭ se li ie reprezentis sian eklezion, kaj fine
li eĉ rifuzadis promesitan subtenon, kiam la libroj ne estis baldaŭ
preparitaj por preso. Kaj tiel ni trovas ĉe Komenský ĵus en tiuj jaroj
1642-48 ofte hezitadon, senkulpiĝadon kaj reĝustiĝadon kaj denove
skuetadon per la katenoj de devoj kiuj lin tiom premis.

15. Cetere, lia maltrankvilo grandiĝis ju pli proksimiĝis la
decido de militaj interkonsentoj, kio signifis ankaŭ decidon pri la
sorto de ĉeĥaj ekzilitoj kaj tamen verdire Komenský sole por Ĉeĥoj
laboris, akceptinte la taskon por li tiom sendankeman. Sed la traktado
estis prokrastata, la esperoj malaperadis kaj ju pli necedema
montriĝis la imperiestro Ferdinando IIIa rilate la demandon de la
religia unueco en bohemaj landoj, des pli volonte permesadis en
Germanio la aprobon de evangeliaj eklezioj. Jam antaŭ la ratifikacio
de pacaj interkonsentoj en la jaro 1648 petas Komenský la reĝinon kaj
svedajn kancelieron, ke ili protektu kompatindajn bohemajn ekzilitojn
kaj ne submetu tiom da senkulpe suferantaj homoj kaj iliaj eklezioj al
la certa ruinigo. Sed ĉio estis vana. Superiĝis la politikaj
konsideroj, ĉar ĉiuj militantaj partioj estis jam lacaj, elĉerpitaj
kaj Ĉeĥoj estis tamen oferdonitaj. Kaj al Komenský ne restis ol
ankoraŭ senkulpigadi sin pro sia insistado.

16. La saman jaron li forlasas Elbingon por administri sian
eklezion, de kiu li estis elektita kiel episkopo--juĝanto. En Lezsno
li poste presigis du librojn, kiujn li dum antaŭaj jaroj verkis por
Svedoj, nome la „Metodon de lingvoj“ kaj „Didaktikon“, sian plej
gravan verkon.

17. Ankaŭ ĉi tie li ne restis longe. Jam la duan jaron post
sia transloĝiĝo en Lesznon li ricevis honoran inviton de princoj
Rákóczi, ke li en Šaryšský Potok super Bodrogo, en loko de ilia
somerloĝado aranĝu lernejon laŭ siaj pansofiaj planoj. Antaŭ sia
forveturo en Hungarion, kien lin sekvis lia eklezio en printempo de la
j. 1650, li ankoraŭ verkis, por ĝojigi la korpremitajn fratojn tie
loĝantajn, „Testamenton de la mortanta patrino de la Frata Unuiĝo,“
verketon de alegoria karaktero, en kiu lia eklezio adiaŭas sur la
mortlito kun siaj filoj kaj postlasas al ili siajn trezorojn. El ĝi
devenas ankaŭ nia enkonduka sentenco, kiu sonas al ni kvazaŭ
profetaĵo. Sed ankaŭ alimaniere ni komprenas la emocion de la
filantropo, kiu ne povas rezigni je la espero en finan venkon de bono,
kvankam la plej proksiman estontecon kovras netraigeblaj nuboj.

18. En Slovakio (Uherská Skalice) li estis je pasko de tiu
jaro tre ĝoje bonveninta; el Uherská Skalice li vizitis ankaŭ aliajn
ekleziojn fratajn pli malproksimen. Iom hezitante, li poste iris
orienton kaj en Potok plenestime akceptis lin precipe la grafo
Zikmundo Rákóczi. Sed por pli longa restado li ekloĝis tie en aŭtuno
de la sama jaro, post longa traktado kun la eklezio, kiu hezitis en
tiu kritika tempo lasi foriri sian ĉefpastron. Fine la politikaj
konsideroj decidis.

19. En sia nova agadloko Komenský tuj prezentis al la princo
planon de la sepklasa lernejo, kies 3 elementaj klasoj estus lingvaj,
4 superaj klasoj objektaj (filozofiaj). La malsuperaj klasoj estis
baldaŭ malfermitaj kaj Komenský relaboris por ili ankaŭ siajn lingvajn
librojn je tria fojo, sed la superaj klasoj ne efektiviĝis pro la tro
frua morto de princo Zikmundo (4./II. 1652.). Komenský volis nome
transmeti la taskon liberigi ĉeĥajn ekzilitojn al Zikmundo, ĉar tiu ĉi
edziĝis kun Eliso, filino de la iama reĝo Frederiko el Pfalz, kiun
vere feliĉan geedziĝon Komenský mem benis. Tial li ankaŭ laboris kun
tia komplezemo por la lernejo, kiun la juna princo tiom
ŝatis. Bedaŭrinde, eĉ tie estis Komenský baldaŭ seniluziita, ĉar la
situacio ne estis ankoraŭ matura por liaj projektoj kaj la rimedoj,
per kiuj li intencis ĉion rebonigi, ne kondukis al la celo. Tial li
fine ellaboris la enhavon de nova „Janua linguarum“ en ok lernejaj
teatraĵoj je „Studo-ludo“, („Schola ludus“), en kiu trovis efektivan
plaĉon nur la nobela junularo de la lernejo. Post la laboro, daŭrinta
3 kaj 1/2 da jaro revenis Komenský alte estimata al sia eklezio en
somero de la jaro 1654.

20. Lian ĝustan kaj praktikan opinion pri la realeco, pruvas
ankaŭ 2 libroj aperintaj en Hungario: Tio estas precipe lia fama
„Orbis pictus“, (Mondo en bildoj), en kiu la konciza enhavo de la
„Janua linguarum“ estas klarigata per multaj bildoj, kiuj estas tie eĉ
la ĉefa temo. La ilustraĵoj anstataŭas tie empirian apercipon, kiun
oni ne povis ĉie akiri, kiel ekzemple en Hungario por diversaj
metioj. Poste, forirante el Hungario, donacis Komenský al la princo
Georgo II., frato de Zikmundo, verketon de ekonomia enhavo, nomitan
„Gentis felicitas“ (Feliĉo de nacio), en kiu Komenský montras, sur kio
bazas la ĝusta signifo de la nacia feliĉo kaj kiel malproksime de ĝi
estas la hungaraj loĝantoj, kaj kiamaniere al ĝi proksimiĝi. Tio certe
estis la plej bela danko por la favoro de Rákóczi.

21. La cirkonstancoj, kiam Komenský la trian fojon alvenis en
Lesznon, estis tre kritikaj, ĉar estis atendita la svedo-pola milito,
kiu ankaŭ vere la sekvantan jaron furioze eksplodis. Pri ĝi ni scias
pli detale el la romano de Sienkiewicz „La inundo.“ La svedan enmarŝon
konvene komparas la poeto kun la inundo, ĉar la malamika militistaro
alproksimiĝis ĝis al Krakow, ke eĉ la reĝo Jan Kazimir mem estis
devigata forkuri fremdlandon. Leszno akceptis la svedan kavalerian
garnizonon. La fremduloj, en la urbo loĝantaj, laŭ religio skismuloj,
estis amataj nek de la popolo el la ĉirkaŭaĵo, nek de la
nobelaro. Kiam do, en unuopaj lokoj de la regno leviĝis ribelo kontraŭ
la fremda regado, venis ankaŭ vico je Leszno, eble ankaŭ la sopiro je
abunda rabaĵo en la riĉa komerca urbo tiam decidis. En aprilo de la
jaro 1656 la kampara popolo, gvidata de la nobelaro, forte atakis la
urbon, kiu post la foriro de la sveda garnizono estis prirabita kaj
ekbruligita. Komenský restis malgraŭ ĉiu minacanta danĝero ĉe sia
paroĥo kaj foriris el la urbo nur en lasta momento, ĵetinte parton de
siaj manuskriptoj en la kaŝejon sub sia studoĉambro. Tie ili estis
ankaŭ trovitaj kaj redonitaj al li dek tagojn post la
brulego. Bedaŭrinde, tio estis sole la malgranda restaĵo de lia
multjara senlaca laboro. Krom tio perdis Komenský sian tutan havaĵon,
la materialon pansofian kaj precipe grandan ĉeĥolatinan vortaron jam
por la preso preparitan. La perdon de tiu ĉi vortara trezoro Komenský
plej multe bedaŭris kaj tio estas ankaŭ la plej granda perdo, kiu
trafis nian nacion je ĝia spirita posedaĵo de la pli malnova
tempo. Pri la valoro de la verko atestas ĝia negranda restaĵo „Linguae
Bohemicae thesaurus“ de V. Jan Rosa, kiun ofte uzadis Josef Jungmann,
verkante sian vortaron.

22. La malfeliĉa sorto de la urbo profunde tuŝis koron precipe
de la samtempa evangelia mondo, same kiel la perdoj de la fama
pansofisto liajn multnombrajn amikojn. Li volis kun sia familio plej
frue transloĝiĝi en Frankfurton a. M., sed volonte li akceptis inviton
de Laurente de Geer, filo de Ludoviko, transloĝiĝi al li en
Amsterodamon kaj tie en trankvilo dediĉi sin al la preferata
pansofio. Loĝantaro de la urbo bonvenigis lin plej estime, kiel
mondfaman personon kaj estis ankaŭ al li proponita helpo por diskonigo
de liaj verkoj.

23. Li akceptis ĝin plej unue por la eldono de la kolektitaj
verkoj pedagogiaj „Opera didactica omnia“ (1657), kiujn li aranĝis laŭ
la tempo de ilia deveno (en Leszno, Elbing kaj en Potok) en tri
grandaj volumoj kaj aldonis al ili kvaran parton de interparoloj,
verkitan en Amsterodamo. Tiuj verkoj ekmirigis la tutan instruitan
mondon; ŝajnis, ke ĉiuj misteroj de la ĝusta instruado estas
malkovritaj kaj nur de tiu unusola viro. La aprobo de la urbo, en kiu
Komenský loĝis kaj al kiu estis tuta verko dediĉita, estis precipe
grava.

24. Kontraŭe la eldono de la latina traduko de politikaj profetaĵoj
devenantaj de tri viziuloj, kun kiuj Komenský persone interrilatis,
vekis kontraŭ li atentindan opozicion kaj maltrankviligis liajn
lastajn jarojn. Kaj tamen ankaŭ per tiu ĉi libro, nomita „Lux in
tenebris“ (Lumo en mallumoj) volis Komenský utili al sia nacio, kaj
denove vaste pritrakti la demandon de ekzilitoj kaj atentigi pri la
maljustaĵoj faritaj al lia nacio. Kion profetis liaj profetoj Kotter
kaj K. Poňatovska pri la unuaj tempoj de la tridek-jara milito estis
ja post la finiĝo de la bataloj plejparte decidite. Krom tio ne ĉesis
instigi Komenský-on lia iama kunlernanto kaj malbonfama profeto el
Lednice Nikolao Drabík al novaj kaj novaj politikaj deklaroj, kvankam
liaj vizioj estis nur senĉese ripetiĝanta atakado je la Habsburga
familio, je Aŭstrio kaj je papismo, de kiuj ruiniĝon en la plej
proksima tempo li profetis. La malfeliĉa rezultato de la militista
interveno en Polio, kiun provis princo Georgo IIa Rákóczi je fino de
la svedo-pola milito, estis certe unu el kaŭzoj de konstanta insistado
de Drabík. Komenský estis tiam ankaŭ peranto de la diplomataj sciigoj
de la princa korto kun la fremdlando. Sian personan honoron kaj sian
neriproĉeblecon li defendis per eldono de la profetaĵoj „Lux e
tenebris“ (Lumo el la malumo, 1665), kiujn li sendis pere de speciala
kuriero al korto de Ludoviko XIV. esperante je komenciĝanta ĵaluzeco
inter Habsburgoj kaj franca reĝa dinastio.

25. Sed malgraŭ la peno por la bono de la subpremitaj
samlandanoj Komenský ĝis la lastaj momentoj ne forgesis, kion li
promesis al la mondo, la pansofion.

El promesitaj „Interkonsiloj universalaj pri la reformo de la
homaj aferoj“ li eldonis nur du enkondukajn volumojn: „Ĝenerala
vekiĝo“--(Panegersio) kaj „Ĝenerala klerigo“--(Panaugiio 1662), sed
li sendis detalan alvokon al ĉiuj eŭropaj scienculoj, volante instigi
klerularon al fervora agado por la leviĝo de la kulturo, sen kiu
ŝajnis al li la homa feliĉo neebla. Krom tio li senĉese primeditadis,
kiamaniere fari la sciencojn alireblajn por ĉiuj klasoj. Tion pruvas
la plano de IIIa parto de la Pansofio konservita bone en lia
postlasaĵo, kiun li fine volis disvolvi en 100 dialogoj. El la pluaj
verkoj ni konas precipe la nomojn: parto IVa „Pampaedia“, „Ĝenerala
kulturo“, kiu estis enhavonta ĉiujn ĉefajn ideojn de la instruado,
parto Va: „Panglottia“; „Universala kulturado de la lingvoj“,
enhavonta krom la ĉefaj lingvaj reguloj ankaŭ ideojn pri la nova
perfekta lingvo por kulturado de la scienco, parto VIa: „Panorthosia,“
„Universala pliboniĝo“, en kiu li volis priskribi la staton de la
homaro, kiam ĉio estos laŭdezire plenumita, parto VIIa: „Pannuthesia,“
„Universala admono“ estis instigonta la homaron, ke ĝi vere pri ĉio ĉi
klopodu. Krom tio Komenský senĉese aliformadis „Trairejon al la
objektoj“ por la junularo kaj li finis la verkon en lasta jaro de sia
vivo, kiel maljunulo preskaŭ okdekjara. Al „Trairejo al la objektoj“
li aldonis kiel praktikan konsekvencon „Triertium catholicum“
(Universala triarto), en kiu li absolute teorie provis solvi la
grandan problemon pri korespondo de la pensado, parolado kaj
agado. Tiun verkon li finskribis komence de la oktobro en 1670, 6
semajnojn antaŭ sia morto.

26. Ankaŭ por la restaĵoj de lia eklezio estas la amsterodama
tempo plena je fruktodona laboro: Komenský ree eldonis tie la ĉefajn
verkojn klerigajn, sed krom tio ankaŭ aliajn, ankoraŭ nepresitajn,
kiel „Manualon“ (mallongigita Sankta Skribo), kantaron, kateĥismon kaj
la konfesion.

27. Kaj kiam li eldonis en la jaro 1668 al la mondo ateston
pri sia vivoagado „Unum necessarium“ (Unu necesa), li povis
kontentiĝinta argumenti, ke li feliĉe trapasis ĉiujn labirintojn, tra
kiuj lin lia kortuŝa vivo kondukis. Li ne sole malgraŭ multaj tentoj
konservis la nefuŝitan kredon de la antauŭloj hereditan, li ne sole
dediĉis grandan parton de sia spirita energio al sia malleviĝanta
eklezio, sed li donis al la tuta mondo kaj precipe al sia nacio brilan
ekzemplon, kiamaniere kaj eĉ dum la plej malĝojaj cirkonstancoj, kiam
ŝajne ĉio la homon konstraŭstaras, oni povas krei la verkon de
senmorta signifo por la tuta estonteco. Pli ol 150 verkojn de Komenský
oni kalkulas, el kiuj ni tie tuŝis nur kelkajn la pli signifajn, sed
ĉiuj pruvas la nerompeblan intelektecon de la menso, senlacan klopodon
klarigi ĉion simple kaj klare. Malmulte de tiel pensantaj kapoj ni
trovas inter la scienculoj de ĉiuj epokoj kaj gentoj.

28. Komenský estas nomata ankaŭ instruisto de la nacioj, ĉar
per sia „Didaktiko“ li montris novajn vojojn de la tuta instruado kaj
krom tio, li mem, persone partoprenis je la nova reformo de la
lernejoj de kelkaj fremdaj nacioj, kiel de Poloj, Germanoj, Svedoj,
Hungaroj kaj Holandanoj. Kaj tamen Komenský mem aldonas, ke ĉion ĉi
tion li volis entrepreni nur por la bono de sia patrujo Bohemio, kiun
li amis per la amo, kiun ne malgrandigis eĉ la longjara ekzilo kaj al
kiu li volis helpi ĝis la lasta spiro de sia vivo.

29. El la ĥaoso de tiama tempo aperas al ni ĝis nun la majesta
staturo de viro de deziro (vir desiderii) staturo de la profeto de la
pli bona estonteco de la homaro, kiam ĉiu ekkonos sian bonon kaj
ekkomprenos, ke la rajtoj ekzistas ne sole por li kaj lia nacio,
provizita per nombra potenco kaj militista forto sed por la tuta
homaro, kaj ke por ĉiu nacio, ĉu granda, ĉu malgranda same valoras
privilegio „Nihil de nobis sine nobis“ (Nenio por ni sen ni). Kaj
Komenský deziris, ke heroldo de tiu ĉi devizo estu ĉiu klerulo, sed
precipe la instruisto. Ankaŭ la paca ligo de la nacioj sub inspekto de
la plej grandaj scienculoj, de veraj amikoj de la homaro, estis ideo
de Komenský. Sed la kulturo ne estu privilegio de iu klaso. Komenský
deziris disvastigi ĝin el ĝia eŭropa lulilo en la plej malproksimajn
landojn de la tero, por ke iam ĉiuj eksentu ĝian benon.`

	return longString[0:bytes]

}
