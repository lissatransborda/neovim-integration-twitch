package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
	"net"
	"github.com/msgpack-rpc/msgpack-rpc-go/rpc"
	"github.com/joho/godotenv"
	"github.com/jrm780/gotirc"
)

func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

func execute(command string) string {
  out, err := exec.Command("bash", "-c", command).Output()

  if err != nil {
    fmt.Println("%s", err)
  }

  output := string(out[:])
  return output
}

func changeVimColor(messageString []string, client *gotirc.Client, channel string, tags map[string]string, clientRPC *rpc.Session){


	if len(messageString) == 2{

		if !strings.ContainsAny(messageString[1], "<>:|") && strings.HasPrefix(messageString[0], "!color"){

			themesAvailble := themes()
			_, found := Find(themesAvailble, messageString[1])
			if found {
				sendInput("<ESC>:color " + messageString[1] + "<CR>", clientRPC)
			}
			if !found{
				client.Say(channel, tags["display-name"] + " Tema não encontrado.")
			}
		}
	}
}

func listThemes(messageString []string, client *gotirc.Client, channel string, tags map[string]string){

	if messageString[0] == "!themes"{

		client.Say(channel, tags["display-name"] + " Veja os temas disponíveis em: https://github.com/edersonferreira/neovim-integration-twitch/blob/main/colors.txt")
	}
}

func move(messageString []string, tags map[string]string, clientRPC *rpc.Session){
	if messageString[0] == "!move" && len(messageString) == 2{
		movement := messageString[1]
		if !strings.ContainsAny(movement, "drRD:<>aioAIOuUvtyYcCsSxX!|-~"){
			sendInput("<ESC>" + messageString[1], clientRPC)
		}
	}
}

func initRPC(port string) *rpc.Session {
	conn, err := net.Dial("tcp", "localhost:" + port)
	if err != nil {
                fmt.Println("fail to connect to server.")
	}
	clientRPC := rpc.NewSession(conn, true)

	return clientRPC
}

func sendInput(input string, clientRPC *rpc.Session){
	_, xerr := clientRPC.Send("nvim_input", input)
	if xerr != nil {
		fmt.Println(xerr)
		return
	}
}

func main() {


  errEnv := godotenv.Load()
  if errEnv != nil {
    log.Fatal("Error loading .env file")
  }

	token := os.Getenv("TOKEN")
	user := os.Getenv("USER")
	channel := os.Getenv("CHANNEL")
	port := os.Getenv("PORT")

	clientRPC := initRPC(port)

        options := gotirc.Options{
            Host:     "irc.chat.twitch.tv",
            Port:     6667,
            Channels: []string{"#" + channel},
        }

        client := gotirc.NewClient(options)
        
        // Whenever someone sends a message, log it
				client.OnChat(func(channel string, tags map[string]string, msg string) {
				
					fmt.Println(msg)
					messageString := strings.Split(msg, " ")

					go listThemes(messageString, client, channel, tags)
					go move(messageString, tags, clientRPC)
					go changeVimColor(messageString, client, channel, tags, clientRPC)
     })

		 client.Connect(user, token)
}

func themes() []string{

	return []string{"0x7A69_dark", "duotone-darkspace", "nightshade","1989", "dusk", "nightshade_print","256-grayvim", "dw_blue", "nightshimmer","256-jungle", "dw_cyan", "nightsky","256_noir", "dw_green", "nightwish","3dglasses", "dw_orange", "no_quarter","Atelier_CaveDark", "dw_purple", "noclown","Atelier_CaveLight", "dw_red", "nocturne","Atelier_DuneDark", "dw_yellow", "nofrils-acme","Atelier_DuneLight", "dzo", "nofrils-dark","Atelier_EstuaryDark", "earendel", "nofrils-light","Atelier_EstuaryLight", "earth", "nofrils-sepia","Atelier_ForestDark", "earthburn", "nord","Atelier_ForestLight", "eclipse", "nordisk","Atelier_HeathDark", "eclm_wombat", "northland","Atelier_HeathLight", "ecostation", "northpole","Atelier_LakesideDark", "editplus", "northsky","Atelier_LakesideLight", "edo_sea", "norwaytoday","Atelier_PlateauDark", "ego", "nour","Atelier_PlateauLight", "eink", "nuvola","Atelier_SavannaDark", "ekinivim", "obsidian","Atelier_SavannaLight", "ekvoli", "obsidian2","Atelier_SeasideDark", "elda", "oceanblack","Atelier_SeasideLight", "eldar", "oceanblack256","Atelier_SulphurpoolDark", "elflord", "oceandeep","Atelier_SulphurpoolLight", "elise", "oceanlight","Benokai", "elisex", "off","Black", "elrodeo", "olive","BlackSea", "elrond", "onedark","Blue2", "emacs", "orange","C64", "enigma", "osx_like","CandyPaper", "enzyme", "otaku","Chasing_Logic", "erez", "oxeded","ChocolateLiquor", "eva", "pablo","ChocolatePapaya", "eva01", "pacific","CodeFactoryv3", "eva01-LCL", "paintbox","Tomorrow", "fx", "pw","Tomorrow-Night", "garden", "py-darcula","Tomorrow-Night-Blue", "gardener", "pyte","Tomorrow-Night-Bright", "gemcolors", "python","Tomorrow-Night-Eighties", "genericdc", "quagmire","VIvid", "genericdc-light", "quantum","White2", "gentooish", "radicalgoodspeed","abbott", "getafe", "raggi","abra", "getfresh", "railscasts","abyss", "ghostbuster", "rainbow_autumn","adam", "github", "rainbow_fine_blue","adaryn", "gobo", "rainbow_fruit","adobe", "golded", "rainbow_night","adrian", "golden", "rainbow_sea","advantage", "goldenrod", "rakr-light","adventurous", "goodwolf", "random","af", "google", "rastafari","afterglow", "gor", "rcg_gui","aiseered", "gotham", "rcg_term","alduin", "gotham256", "rdark","ancient", "gothic", "rdark-terminal","anderson", "grape", "redblack","angr", "gravity", "redstring","anotherdark", "grayorange", "refactor","ansi_blows", "graywh", "relaxedgreen","antares", "grb256", "reliable","apprentice", "greens", "reloaded","aqua", "greenvision", "revolutions","aquamarine", "greenwint", "robinhood","arcadia", "grey2", "rockets-away","archery", "greyblue", "ron","argonaut", "greygull", "rootwater","ashen", "grishin", "sadek1","asmanian2", "gruvbox", "sand","Dark", "evening", "paramount","Dark2", "evening1", "parsec","DarkDefault", "evokai", "peachpuff","DevC++", "evolution", "peaksea","Dev_Delight", "fahrenheit", "pencil","Dim", "fairyfloss", "penultimate","Dim2", "falcon", "peppers","DimBlue", "far", "perfect","DimGreen", "felipec", "petrel", "DimGreens", "feral", "pf_earth","DimGrey", "fight-in-the-shade", "phd","DimRed", "fine_blue", "phoenix","DimSlate", "firewatch", "phphaxor","Green", "flatcolor", "phpx","Light", "flatland", "pink","LightDefault", "flatlandia", "pixelmuerto","LightDefaultGrey", "flattened_dark", "plasticine","LightTan", "flattened_light", "playroom","LightYellow", "flattown", "pleasant","Monokai", "flattr", "potts","MountainDew", "flatui", "predawn","OceanicNext", "fnaqevan", "preto","OceanicNextLight", "fog", "pride","PapayaWhip", "fokus", "primaries","PaperColor", "forneus", "primary","PerfectDark", "foursee", "print_bw","Red", "freya", "prmths","Revolution", "frood", "professional","SerialExperimentsLain", "frozen", "proton","Slate", "fruidle", "ps_color","SlateDark", "fruit", "pspad","Spink", "fruity", "pt_black","SweetCandy", "fu", "putty","asmanian_blood", "gryffin", "sandydune","asmdev", "guardian", "satori","asmdev2", "guepardo", "saturn","astronaut", "h80", "scheakur","asu1dark", "habiLight", "scite","atom", "happy_hacking", "scooby","aurora", "harlequin", "seagull","automation", "heliotrope", "sean","autumn", "hemisu", "seashell","autumnleaf", "herald", "seattle","ayu", "heroku", "selenitic","babymate256", "heroku-terminal", "seoul","badwolf", "herokudoc", "seoul256","bandit", "herokudoc-gvim", "seoul256-light","base", "hhazure", "seti","base16-ateliercave", "hhdblue", "settlemyer","base16-atelierdune", "hhdcyan", "sexy-railscasts","base16-atelierestuary", "hhdgray", "sf","base16-atelierforest", "hhdgreen", "shades-of-teal","base16-atelierheath", "hhdmagenta", "shadesofamber","base16-atelierlakeside", "hhdred", "shine","base16-atelierplateau", "hhdyellow", "shiny-white","base16-ateliersavanna", "hhorange", "shobogenzo","base16-atelierseaside", "hhpink", "sialoquent","base16-ateliersulphurpool", "hhspring", "sienna","base16-railscasts", "hhteal", "sierra","basic", "hhviolet", "sift","basic-dark", "highlighter_term", "silent","basic-light", "highlighter_term_bright", "simple256","bayQua", "highwayman", "simple_b","baycomb", "hilal", "simple_dark","bclear", "holokai", "simpleandfriendly","beachcomber", "hornet", "simplewhite","beauty256", "horseradish256", "simplon","beekai", "hotpot", "skittles_autumn","behelit", "hual", "skittles_berry","benlight", "hybrid", "skittles_dark","bensday", "hybrid-light", "sky","billw", "hybrid_material", "slate","biogoo", "hybrid_reverse", "slate2","birds-of-paradise", "hydrangea", "smarties","bitterjug", "iangenzo", "smp","black_angus", "ibmedit", "smpl","blackbeauty", "icansee", "smyck","blackboard", "iceberg", "soda","blackdust", "immortals", "softblue","blacklight", "impact", "softbluev2","blaquemagick", "impactG", "softlight","blazer", "impactjs", "sol","blink", "industrial", "sol-term","blue", "industry", "solarized","bluechia", "ingretu", "solarized8_dark","bluedrake", "inkpot", "solarized8_dark_flat","bluegreen", "inori", "solarized8_dark_high","bluenes", "ir_black", "solarized8_dark_low","blueprint", "ironman", "solarized8_light","blues", "itg_flat", "solarized8_light_flat","blueshift", "itg_flat_transparent", "solarized8_light_high","bluez", "itsasoa", "solarized8_light_low","blugrine", "jaime", "sole","bluish", "jammy", "sonofobsidian","bmichaelsen", "janah", "sonoma","boa", "japanesque", "sorcerer","bocau", "jelleybeans", "soruby","bog", "jellybeans", "soso","boltzmann", "jellygrass", "sourcerer","borland", "jellyx", "southernlights","breeze", "jhdark", "southwest-fog","cascadia", "leya", "tango","celtics_away", "lightcolors", "tango-desert","cgpro", "lightning", "tango-morning","chalkboard", "lilac", "tango2","chance-of-storm", "lilydjwg_dark", "tangoX","charged-256", "lilydjwg_green", "tangoshady","charon", "lilypink", "taqua","chela_light", "lingodirector", "tatami","cherryblossom", "liquidcarbon", "tayra","chlordane", "literal_tango", "tchaba","chocolate", "lizard", "tchaba2","chroma", "lizard256", "tcsoft","chrysoprase", "lodestone", "telstar","clarity", "loogica", "tender","cleanphp", "louver", "termschool","cleanroom", "lucid", "tesla","clearance", "lucius", "tetragrammaton","cloudy", "luinnar", "textmate16","clue", "lumberjack", "thegoodluck","cobalt", "luna", "thermopylae","cobalt2", "luna-term", "thestars","cobaltish", "lxvc", "thor","coda", "lyla", "thornbird","codeblocks_dark", "mac_classic", "tibet","codeburn", "macvim-light", "tidy","codedark", "made_of_code", "tigrana-256-dark","codeschool", "madeofcode", "tigrana-256-light","coffee", "magellan", "tir_black","coldgreen", "magicwb", "tolerable","colorer", "mango", "tomatosoup","colorful", "manuscript", "tony_light","colorful256", "manxome", "toothpik","colorsbox-faff", "marklar", "torte","colorsbox-greenish", "maroloccio", "transparent","breezy", "jhlight", "space-vim-dark","brighton", "jiks", "spacegray","briofita", "jitterbug", "spacemacs-theme","broduo", "kalahari", "spartan","brogrammer", "kalisi", "spectro","brookstream", "kalt", "spiderhawk","brown", "kaltex", "spring","bubblegum", "kate", "spring-night","bubblegum-256-dark", "kellys", "sprinkles","bubblegum-256-light", "khaki", "spurs_away","buddy", "kib_darktango", "srcery","burnttoast256", "kib_plastic", "srcery-drk","busierbee", "kings-away", "stackoverflow","busybee", "kiss", "stefan","buttercream", "kkruby", "stereokai","bvemu", "koehler", "stingray","bw", "kolor", "stonewashed-256","c", "kruby", "stonewashed-dark-256","c16gui", "kyle", "stonewashed-dark-gui","cabin", "laederon", "stonewashed-gui","cake", "lakers_away", "stormpetrel","cake16", "landscape", "strange","calmar256-dark", "lanox", "strawimodo","calmar256-light", "lanzarotta", "summerfruit","camo", "lapis256", "summerfruit256","campfire", "last256", "sunburst","candy", "late_evening", "surveyor","candycode", "lazarus", "swamplight","candyman", "legiblelight", "sweater","caramel", "leglight2", "symfony","carrot", "leo", "synic","carvedwood", "less", "synthwave","carvedwoodcool", "lettuce", "tabula","colorsbox-material", "maroloccio2", "triplejelly","colorsbox-stblue", "maroloccio3", "trivial256","colorsbox-stbright", "mars", "trogdor","colorsbox-steighties", "martin_krischik", "tropikos","colorsbox-stnight", "material", "true-monochrome","colorzone", "material-theme", "turbo","contrastneed", "materialbox", "turtles","contrasty", "materialtheme", "tutticolori","cool", "matrix", "twilight","corn", "maui", "twilight256","corporation", "mayansmoke", "twitchy","crayon", "mdark", "two-firewatch","crt", "mellow", "two2tango","crunchbang", "messy", "ubaryd","cthulhian", "meta5", "ubloh","custom", "metacosm", "umber-green","cyberpunk", "midnight", "understated","d8g_01", "miko", "underwater","d8g_02", "minimal", "underwater-mod","d8g_03", "minimalist", "unicon","d8g_04", "mint", "up","dalton", "mizore", "valloric","dante", "mod8", "vanzan_color","dark-ruby", "mod_tcsoft", "vc","darkBlue", "mohammad", "vcbc","darkZ", "mojave", "vertLaiton","darkblack", "molokai", "vexorian","darkblue", "molokai_dark", "vibrantink","darkblue2", "monoacc", "vice","darkbone", "monochrome", "vilight","darkburn", "monokai-chris", "vim-material","darkdevel", "monokai-phoenix", "vimbrains","darkdot", "monokain", "vimbrant","darkeclipse", "montz", "vimicks","darker-robin", "moody", "visualstudio","darkerdesert", "moonshine", "vividchalk","darkglass", "moonshine_lowcontrast", "vj","darkocean", "moonshine_minimal", "void","darkrobot", "mophiaDark", "vorange","darkslategray", "mophiaSmoke", "vydark","darkspectrum", "mopkai", "vylight","darktango", "more", "wargrey","darkzen", "moria", "warm_grey","darth", "moriarty", "warriors-away","dawn", "morning", "wasabi256","deep-space", "moss", "watermark","deepsea", "motus", "wellsokai","default", "mourning", "welpe","delek", "mrkn256", "white","delphi", "mrpink", "whitebox","denim", "mud", "whitedust","derefined", "muon", "widower","desert", "murphy", "wikipedia","desert256", "mushroom", "win9xblueback","desert256v2", "mustang", "winter","desertEx", "mythos", "winterd","desertedocean", "native", "wintersday","desertedoceanburnt", "nature", "woju","desertink", "navajo", "wolfpack","despacio", "navajo-night", "wombat","detailed", "nazca", "wombat256","deus", "nedit", "wombat256dave","devbox-dark-256", "nedit2", "wombat256i","deveiate", "nefertiti", "wombat256mod","developer", "neodark", "wood","diokai", "neon", "wuye","disciple", "neonwave", "wwdc16","distill", "nerv-ous", "wwdc17","distinguished", "nes", "xcode","django", "nets-away", "xcode-default","donbass", "neuromancer", "xedit","donttouchme", "neutron", "xemacs","doorhinge", "neverland", "xian","doriath", "neverland-darker", "xmaslights","dual", "neverland2", "xoria256","dull", "neverland2-darker", "xterm16","duotone-dark", "neverness", "yeller","duotone-darkcave", "nevfn", "yuejiu","duotone-darkdesert", "new-railscasts", "zazen","duotone-darkearth", "newspaper", "zellner","duotone-darkforest", "newsprint", "zen","duotone-darkheath", "nicotine", "zenburn","duotone-darklake", "night", "zenesque","duotone-darkmeadow", "nightVision", "zephyr","duotone-darkpark", "night_vision", "zmrok","duotone-darkpool", "nightflight", "znake","duotone-darksea", "nightflight2"}

}
