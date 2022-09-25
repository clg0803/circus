package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/clg0803/circus/evaluator"
	"github.com/clg0803/circus/lexer"
	"github.com/clg0803/circus/object"
	"github.com/clg0803/circus/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvirnment()
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		eval := evaluator.Eval(program, env)
		if eval != nil {
			io.WriteString(out, eval.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
const IKUN = `
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,:,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~,,,,,,,,,,,,,,,,,,,,,,,,IMMMMMMMMM=,,,,,,,,,,,,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~,,,,,,,,,,,,,,,,,,,,,?MMMMMNMMMMMMMMM:,,,,,,,,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~:,,,,,,$MMMMMMMD=,,,MMMM??======+ZMMMMMN:,,,,,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~,:7MMMMMMMMMMMMMMMMMMZ==,~=~..,,:==+MMMMMZ,,,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~DMMMMMMI+=======ID+NMI~:.:::::::~==:===IMMMMN,,,,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~NMMMMM7+======,,,,,,==,,::,::~:~:::::::~=+==IDMMMM:,,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~IMMMMMI======+I?=~:::::~~::::,,=======::::~:~~==IIDMMMN,,,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~NMMMM$?I+==?IIIII?::::,,,==MMM+,=NMMMN====,,:::::~==IIDMMMD,,,,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~~~NMMMMII~,=?III??:::::::::~=+MMM==NMMMMMMMD=====::::,.:+=?IMMMM=:,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~~~MMMMDII,,=?III=::::===::::~==MMMIOMMMNIIIMMMMZI===,~::::,===+IMMMM,,,,,,,,,
~~~~~~~~~~~~~~~~~~~~~~~8MMMNI?.,=IIII~::,::=====:::==ZMMZZMMMMMMMMMMMMMMI+==,::::~,,===IZMMM=,,,,,,,
~~~~~~~~~~~~~~~~~~~~~?MMMNI+,,+?II?::::.~II====,::===MMMMMMMMMMMMMMMMMMMM8?==,8MM:::,===?IMMMN,,,,,,
~~~~~~~~~~~~~~~~~~~~MMMMII:,==?II~:~::,~II?===,:+~===MMMMM7?IIIIIIIIIIIMMMNI==?MMM~~:,~+=?IMMMD:,,,,
~~~~~~~~~~~~~~~~~~=MMMO?=,====II~::::.~II?===::=II==OMMIIIIII~~~~~+IIIIINMMDI===MMM$::,:==??NMMZ,,,,
,~~~~~~~~~~~~~~~~ZMMM??.,====:~~==~:,,III===~::II?==MMM?IIII~~~~~~~~?IIIIMMMII==IMMM?~:,:==IIMMM,,,,
,,,:~~~~~~~~~~~:DMMMI+,+===,:?MMO+=:,+II====:::II+==MMMMMN8~:~~~~~~~~~IIIIMMMI===IMMM=::,~==I?MMM,,,
,,,,,:~~~~~~~~~MMMNI~.====.:~MMM+==,:II+===:,:=II===MMMMMMMMMMO~~~~~~~~+IIMMMZI,==?MMM~::,==+IOMM?,,
,,,,,,,:~~~~~~NMMOI~,====:,:ZMMI==:,?II===~,::II?==+MMO==,..DNMMN~~~~~~?MMMMMMI::==IMMM::~,==I?MMM,,
,,,,,,,,,,~~~OMMNI=,==?I+.:~MMM===.~II+===,:::II~==?MM7==.....=MMM~~~~MMMZ=MMMII.==?ZMMN:::~==I8MM~,
,,,,,,,,,,,,~MMMI+MMM7II=,::MMN+=~,+II+==~.::=II===7MM?==.......MMM+:MMM.,+OMMMM7==+?NMM?~::====MMN,
,,,,,,,,,,,,NMMI+=MMM+II,,:?MM7==,,II?===,:~:======$MM+=~........MMMNMM,..=IMMMMMMMMMNMMM~::====MMM:
,,,,,,,,,,,,MMM7=IMMI=II,::DMMI==,:II====.:::=NI===$MM?=,..:MMMN..MMMM....=====+MMMMMMMMMM======NMM,
,,,,,,,,,,,,,MMMMNMM=+II,~:MMMI==,~II===?,::=NMMMN=+MM7=,..MMMMMM.MMMM.....====MMMMMM++MMM==?I?=8MM:
,,,,,,,,,,,,,,+MMMMM=+II:~:MMMI==,:II==MMM~===MMMMMMMMN=..~MMMMMM.DMMM........:MMMMMM+=MMMM=?II=MMM,
,,,,,,,,,,,,,,,,:MMM+?II~::MMMI=~,:II==MMN=+MMMM8IMMMMM=~..NMMMM$.MMMM.........NMMMMN=MMMMMD+IIMMM=,
,,,,,,,,,,,,,,,,,MMM=+II:==MMMI+=~:??=IMMMMMMMMMM====+===........:MMMMO......... .,..MMMINMMMMMMMO,,
,,,,,,,,,,,,,,,,,MMM=+II===NMMI?=::~==ZMMMMD?IIMMM:.,~===.......,MMM=MM$............MMMIII?MMMMM?,,,
,,,,,,,,,,,,,,,,,MMM==I7MMMMMM7I==,===ZMMZOIIII:MMMM,..........MMMN~=~MMM,........MMMN~=7ZO8MM~,,,,,
,,,,,,,,,,,,,,,,,OMMMMMMMMMMMMMI?=====+MMDZO7~~~~:MMMMM~...=NMMMMMMMMMMMMMMMD88MMMMM~=$$777$MM8,,,,,
,,,,,,,,,,,,,,,,,,NMMMMMM777NMMMMMMMMMMMMNZZ$?~~~~~=7MMMMMMMMMMI+++++++++8MMMMMMM+~~+$77777$NMM,,,,,
,,,,,,,,,,,,,,,,,,,,,:=MMI+77MMMMMMMMMNOZOO$7$~~~~~~~~~~~~MMM++++++++++++++?MMD~~~~~$$777777NMM,,,,,
,,,,,,,,,,,::,~8MMMMMMMMM$==ZZOOZZOZOOZOZ$777$=~~~~~~~~~~MMN+++++++++++++++++MM?~~~I$7777777NMM:,,,,
,,,,,,,,,,?MMMMMMMMDZZZMMM==IZZOOZZ$77$77$777$=~~~~~~~~~~MM?+++++++++++++++IMMMM~~~$$7777777MMM,,,,,
,,,,,,,IMMMMM7II77II7$$DMM=++7777777777777777$~~~~~~~~~~~MMMMD==++++=?+7MMMMMMMD~~~$7777777$MM7,,,,,
,,,,:OMMMDII7I777IIII7$$MMM==7$7777777777777$==~======~~~DMMMMMMMMMMMMMMMM8+ZMM~~~~?$777777NMM:,,,,,
,,,IMMMMMMMMMMDZ7777Z8NMMMMD++I$7$777777777$+=======++++==$MMM+++++++++++++NMM~~~~~~$77777DMM~,,,,,,
,,MMMZII7IOMMMMMMMMMMMMMMMMMD=++$$7777777$I+++=======++++++=8MMMN?++++++IMMMO~~~~~~~=$7$7NMM7,,,,,,,
,MMM7I7II7I77777IIIIIIII$$DMMM==+==?77I?+====++====+=++++=+====?MMMMMMMMMD===~~~~~~~~~7ZMMM?~:,,,,,,
MMMII7II777777777777777I$$MMMMMN==+=====ONMMMMMMMMMMMMMMMMMN$=====+=====+=+=+=+++===~7MMMM:~~~~:,,,,
MMMMMMMMN7IIIIIII7777777$ZMMMMMMMMMMMMMMMMMMMMMMMNMMMMMMMMMMMMMMMMMMN8+===========NMMMMMMMMD~~~~~~,,
MII778MMMMMMMM7II7777777$8MMMMMMMMMMNNNDDNDNNNDNNNNDNDDDNNNNDNDNNMMMMMMMMMMMN7====?MMNDNNMMMMMN~~~~~
NI7I7777I77NMMMMMMMZIIII$OMMMMMNDNNNNNNNNNNNNNNNNNNNNDDDNNNDNNNNDDDNNNNNNNMMMMMMMMMMMMMMMMMMMMO~~~~~
I7I7I7III7II77I7NMMMMMMNZ$MMMDNDDNNNNNNNNNNNNNNNNNDNDNNDNMMMMMMMMMMMNNNDNDNDDDDNMMMMMMMMMMMMMM~~~~~~
7IIII777IIIII7I7III7IDMMMMMMMMNNDNNNNNNNNNNNNNNNNNDNMMMMMMMMMMMMMMMMMMMMMMNDDNNDDDNNNNNNNNMMM~~~~~~~
7I7I77I7IIIIIII7IIIIIIII$NMMMMMMMMMMMMMNNNNNNMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNNNNNMMMMMDNM~~~~~~~
I7II7777IIIIIIIIIIIIIIII$M+~:::::::::~ZMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM=:::::?M~~~~~~~
MMMMMMN$IIIIIII7IIIIII77$M?=====~=~::::::MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM~:::~====+M~~~~~~~
7ODNMMMMMMMMNII7III7III7$$MMMMMMMMN===~:::~MMNDNNNDNNMMMMMMMMMMMMMMMMMMMMMMMMMMM=::~===8MMMMD~~~~~~~
77IIIIIII$MMMMMM7777III7$$$MMNNNMMMMMM==:::$M+~~=$DNNNNNNDNNNNNNDNNNNNNZ=~~==OM:::==$MMMMNMMM~~~~~~~
MI77IIIIII7II7NMMMNI7I77$$MMMDNNNNNMMMMM=~::MNNND8+=~=:~+8DDDDDO?======8MDNNNM~::=?MMMNNNNNMMN~~~~~~
MIIIIIIIII7777I7IMMMMI7$$$MMNDNNNNNNNMMM$=::NMNNNNNNDNDDDI=====?NMNNNNDNNNNNMM::=?MMMNNNNNNNMM8~~~~~
MMI7I777II77IIII777NMMM$$$MMDDDNNNNNDDMMM=::OMNNNNNNNNNNNNDD?=DMDDNNNNNNNNNNMN::=MMMNNNNNNNNNMM~~~~~
MMM77777II77II777777INMMNOMMNNNNNNNNNNMMM=::OMDNNNNNNNNNNN~===~=+NNNNNNNNNNNMI:~=MMMNDNNNNNNNMMM~~~~
,MMMI777IIIIIIIIIIIIII7MMMMMNNNNNNNNNNMMM=::DMDNNNNNNNNNNN~=~~::=DNNNNNNNNNNM::~=MMNNDNNNNNNNNMMN~~~
,,MMM777IIIIIIIIIIIIII77DMMMDDNNNNNNNNMMM=::MMDNNNNNNNNNND~~,=~=~?NNNNNNNNNNM=:==MMMNNNNNNNNNNNMM=~~

`

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, IKUN)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
