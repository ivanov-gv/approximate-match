package integration

// StationData holds all information about one station (or blacklisted location)
// in a single unified record.
//
// Official stations have a positive Id matching the ZPCG station database.
// Blacklisted locations have a negative Id and Blacklisted=true; they exist in
// the dataset so tests can verify the matcher does NOT surface them.
type StationData struct {
	Id         int
	ZpcgStopId int
	Type       int

	// Official names — indexed by the official matcher.
	Name    string // primary Serbian Latin name
	NameEn  string // English name (may equal Name)
	NameCyr string // Serbian Cyrillic name

	// ProductionAliases are additional names used alongside the official ones in
	// the production system to prevent disambiguation (e.g. short forms, Russian
	// Cyrillic spellings, English translations).
	ProductionAliases []string

	// UserInputVariants are phonetic or transliteration forms that real users
	// type but that are not production names (e.g. folk-etymology spellings,
	// Russian variants that diverge too far to be caught by fuzzy matching alone).
	UserInputVariants []string

	// Blacklisted marks locations that exist near the railway network but have
	// no active service. The matcher should NOT return them when the official
	// station list is searched.
	Blacklisted bool
}

// Stations is the single source of truth for all station test data.
// Official stations (Id > 0) are followed by blacklisted locations (Id < 0).
var Stations = []StationData{
	{
		Id: 1, ZpcgStopId: 1, Type: 4,
		Name: "Bar", NameEn: "Bar", NameCyr: "Бар",
	},
	{
		Id: 2, ZpcgStopId: 3, Type: 1,
		Name: "Sutomore", NameEn: "Sutomore", NameCyr: "Сутоморе",
	},
	{
		Id: 3, ZpcgStopId: 6, Type: 1,
		Name: "Golubovci", NameEn: "Golubovci", NameCyr: "Голубовци",
		ProductionAliases: []string{"Голубовцы"},
	},
	{
		Id: 4, ZpcgStopId: 8, Type: 4,
		Name: "Podgorica", NameEn: "Podgorica", NameCyr: "Подгорица",
	},
	{
		Id: 5, ZpcgStopId: 15, Type: 1,
		Name: "Kolašin", NameEn: "Kolašin", NameCyr: "Колашин",
	},
	{
		Id: 6, ZpcgStopId: 18, Type: 1,
		Name: "Mojkovac", NameEn: "Mojkovac", NameCyr: "Мојковац",
		ProductionAliases: []string{"Мойковац"},
	},
	{
		Id: 7, ZpcgStopId: 22, Type: 4,
		Name: "Bijelo Polje", NameEn: "Bijelo Polje", NameCyr: "Бијело Поље",
		ProductionAliases: []string{"Бело Поле"},
	},
	{
		Id: 8, ZpcgStopId: 55, Type: 1,
		Name: "Prijepolje teretna", NameEn: "Prijepolje cargo", NameCyr: "Пријепоље теретна",
		ProductionAliases: []string{"Приеполье теретна"},
	},
	{
		Id: 9, ZpcgStopId: 26, Type: 2,
		Name: "Prijepolje", NameEn: "Prijepolje", NameCyr: "Пријепоље",
		ProductionAliases: []string{"Приеполье"},
	},
	{
		Id: 10, ZpcgStopId: 27, Type: 2,
		Name: "Priboj", NameEn: "Priboj", NameCyr: "Прибој",
		ProductionAliases: []string{"Прибой"},
	},
	{
		Id: 11, ZpcgStopId: 28, Type: 2,
		Name: "Užice", NameEn: "Užice", NameCyr: "Ужице",
	},
	{
		Id: 12, ZpcgStopId: 29, Type: 3,
		Name: "Požega", NameEn: "Požega", NameCyr: "Пожега",
	},
	{
		Id: 13, ZpcgStopId: 62, Type: 1,
		Name: "Kosjerić", NameEn: "Kosjeric", NameCyr: "Косјерић",
		ProductionAliases: []string{"Косерич"},
	},
	{
		Id: 14, ZpcgStopId: 30, Type: 2,
		Name: "Valjevo", NameEn: "Valjevo", NameCyr: "Ваљево",
		ProductionAliases: []string{"Вальево"},
	},
	{
		Id: 15, ZpcgStopId: 66, Type: 1,
		Name: "Lajkovac", NameEn: "Lajkovac", NameCyr: "Лајковац",
		ProductionAliases: []string{"Лайковац"},
	},
	{
		Id: 16, ZpcgStopId: 65, Type: 1,
		Name: "Lazarevac", NameEn: "Lazarevac", NameCyr: "Лазаревац",
	},
	{
		Id: 17, ZpcgStopId: 61, Type: 1,
		Name: "Rakovica", NameEn: "Rakovica", NameCyr: "Раковица",
	},
	{
		Id: 18, ZpcgStopId: 31, Type: 4,
		Name: "Beograd Centar", NameEn: "Belgrade Center", NameCyr: "Београд Центар",
		ProductionAliases: []string{"Beograd", "belgrad", "belgrade", "Белград Центр"},
	},
	{
		Id: 19, ZpcgStopId: 75, Type: 1,
		Name: "Zemun", NameEn: "Zemun", NameCyr: "Земун",
	},
	{
		Id: 20, ZpcgStopId: 2, Type: 2,
		Name: "Šušanj", NameEn: "Šušanj", NameCyr: "Шушањ",
		ProductionAliases: []string{"Шушань"},
	},
	{
		Id: 21, ZpcgStopId: 43, Type: 2,
		Name: "Crmnica", NameEn: "Crmnica", NameCyr: "Црмница",
	},
	{
		Id: 22, ZpcgStopId: 4, Type: 1,
		Name: "Virpazar", NameEn: "Virpazar", NameCyr: "Вирпазар",
	},
	{
		Id: 23, ZpcgStopId: 44, Type: 2,
		Name: "Vranjina", NameEn: "Vranjina", NameCyr: "Врањина",
		ProductionAliases: []string{"Враньина"},
	},
	{
		Id: 24, ZpcgStopId: 5, Type: 3,
		Name: "Zeta", NameEn: "Zeta", NameCyr: "Зета",
	},
	{
		Id: 25, ZpcgStopId: 45, Type: 2,
		Name: "Morača", NameEn: "Morača", NameCyr: "Морача",
	},
	{
		Id: 26, ZpcgStopId: 7, Type: 2,
		Name: "Aerodrom", NameEn: "Aerodrom", NameCyr: "Аеродром",
		ProductionAliases: []string{"аэродром"},
	},
	{
		Id: 27, ZpcgStopId: 9, Type: 2,
		Name: "Zlatica", NameEn: "Zlatica", NameCyr: "Златица",
	},
	{
		Id: 28, ZpcgStopId: 46, Type: 3,
		Name: "Bioče", NameEn: "Bioče", NameCyr: "Биоче",
	},
	{
		Id: 29, ZpcgStopId: 10, Type: 3,
		Name: "Bratonožići", NameEn: "Bratonožići", NameCyr: "Братоножићи",
		ProductionAliases: []string{"Братоношићи", "Братоношичи"},
	},
	{
		Id: 30, ZpcgStopId: 11, Type: 3,
		Name: "Lutovo", NameEn: "Lutovo", NameCyr: "Лутово",
	},
	{
		Id: 31, ZpcgStopId: 47, Type: 2,
		Name: "Kruševački Potok", NameEn: "Kruševački Potok", NameCyr: "Крушевачки Поток",
		ProductionAliases: []string{"Крушевачки поток", "Крушевацки поток"},
	},
	{
		Id: 32, ZpcgStopId: 12, Type: 1,
		Name: "Trebešica", NameEn: "Trebešica", NameCyr: "Требешица",
	},
	{
		Id: 33, ZpcgStopId: 13, Type: 2,
		Name: "Selište", NameEn: "Selište", NameCyr: "Селиште",
	},
	{
		Id: 34, ZpcgStopId: 48, Type: 3,
		Name: "Kos", NameEn: "Kos", NameCyr: "Кос",
	},
	{
		Id: 35, ZpcgStopId: 14, Type: 2,
		Name: "Mateševo", NameEn: "Mateševo", NameCyr: "Матешево",
		ProductionAliases: []string{"Матесево"},
	},
	{
		Id: 36, ZpcgStopId: 49, Type: 2,
		Name: "Padež", NameEn: "Padež", NameCyr: "Падеж",
	},
	{
		Id: 37, ZpcgStopId: 16, Type: 2,
		Name: "Oblutak", NameEn: "Oblutak", NameCyr: "Облутак",
	},
	{
		Id: 38, ZpcgStopId: 50, Type: 3,
		Name: "Trebaljevo", NameEn: "Trebaljevo", NameCyr: "Требаљево",
		ProductionAliases: []string{"Требальево"},
	},
	{
		Id: 39, ZpcgStopId: 17, Type: 2,
		Name: "Štitarička Rijeka", NameEn: "Štitarička Rijeka", NameCyr: "Штитаричка Ријека",
		ProductionAliases: []string{"Штитарица река", "Штитаричка река"},
	},
	{
		Id: 40, ZpcgStopId: 51, Type: 2,
		Name: "Žari", NameEn: "Žari", NameCyr: "Жари",
		ProductionAliases: []string{"Зари"},
	},
	{
		Id: 41, ZpcgStopId: 19, Type: 3,
		Name: "Mijatovo Kolo", NameEn: "Mijatovo Kolo", NameCyr: "Мијатово Коло",
		ProductionAliases: []string{"Мијатово коло", "Миятово коло"},
	},
	{
		Id: 42, ZpcgStopId: 52, Type: 2,
		Name: "Slijepač Most", NameEn: "Slijepač Most", NameCyr: "Слијепач Мост",
		ProductionAliases: []string{"Слијепац мост", "Слепец мост"},
	},
	{
		Id: 43, ZpcgStopId: 53, Type: 2,
		Name: "Ravna Rijeka", NameEn: "Ravna Rijeka", NameCyr: "Равна Ријека",
		ProductionAliases: []string{"Равна ријека", "Равна река"},
	},
	{
		Id: 44, ZpcgStopId: 20, Type: 3,
		Name: "Kruševo", NameEn: "Kruševo", NameCyr: "Крушево",
	},
	{
		Id: 45, ZpcgStopId: 21, Type: 2,
		Name: "Lješnica", NameEn: "Lješnica", NameCyr: "Љешница",
		ProductionAliases: []string{"Льешница"},
	},
	{
		Id: 46, ZpcgStopId: 35, Type: 2,
		Name: "Pričelje", NameEn: "Pričelje", NameCyr: "Причеље",
		ProductionAliases: []string{"Прицеље", "Причелье"},
	},
	{
		Id: 47, ZpcgStopId: 23, Type: 1,
		Name: "Spuž", NameEn: "Spuž", NameCyr: "Спуж",
	},
	{
		Id: 48, ZpcgStopId: 36, Type: 2,
		Name: "Ljutotuk", NameEn: "Ljutotuk", NameCyr: "Љутотук",
		ProductionAliases: []string{"Лютотук"},
	},
	{
		Id: 49, ZpcgStopId: 24, Type: 4,
		Name: "Danilovgrad", NameEn: "Danilovgrad", NameCyr: "Даниловград",
	},
	{
		Id: 50, ZpcgStopId: 37, Type: 2,
		Name: "Slap", NameEn: "Slap", NameCyr: "Слап",
	},
	{
		Id: 51, ZpcgStopId: 38, Type: 2,
		Name: "Bare Šumanovića", NameEn: "Bare Šumanovića", NameCyr: "Баре Шумановића",
		ProductionAliases: []string{"Баре Шумановица"},
	},
	{
		Id: 52, ZpcgStopId: 39, Type: 2,
		Name: "Šobajići", NameEn: "Šobajići", NameCyr: "Шобајићи",
		ProductionAliases: []string{"Собајићи", "Шобайичи"},
	},
	{
		Id: 53, ZpcgStopId: 40, Type: 3,
		Name: "Ostrog", NameEn: "Ostrog", NameCyr: "Острог",
	},
	{
		Id: 54, ZpcgStopId: 41, Type: 2,
		Name: "Dabovići", NameEn: "Dabovići", NameCyr: "Дабовићи",
		ProductionAliases: []string{"Дабовичи"},
	},
	{
		Id: 55, ZpcgStopId: 42, Type: 2,
		Name: "Stubica", NameEn: "Stubica", NameCyr: "Стубица",
	},
	{
		Id: 56, ZpcgStopId: 25, Type: 4,
		Name: "Nikšić", NameEn: "Nikšić", NameCyr: "Никшић",
		ProductionAliases: []string{"Никшич"},
	},
	{
		Id: 57, ZpcgStopId: 32, Type: 4,
		Name: "Novi Sad", NameEn: "Novi Sad", NameCyr: "Нови Сад",
		ProductionAliases: []string{"Нови сад", "Новый Сад"},
	},
	{
		Id: 58, ZpcgStopId: 33, Type: 2,
		Name: "Vrbas", NameEn: "Vrbas", NameCyr: "Врбас",
	},
	{
		Id: 59, ZpcgStopId: 34, Type: 4,
		Name: "Subotica", NameEn: "Subotica", NameCyr: "Суботица",
	},
	{
		Id: 60, ZpcgStopId: 54, Type: 1,
		Name: "Novi Beograd", NameEn: "Novi Beograd", NameCyr: "Нови Београд",
		ProductionAliases: []string{"Novi belgrad", "New Belgrade", "Новый Белград"},
	},
	{
		Id: 61, ZpcgStopId: 56, Type: 1,
		Name: "Čačak", NameEn: "Čačak", NameCyr: "Чачак",
	},
	{
		Id: 62, ZpcgStopId: 57, Type: 1,
		Name: "Kraljevo", NameEn: "Kraljevo", NameCyr: "Краљево",
	},
	{
		Id: 63, ZpcgStopId: 58, Type: 1,
		Name: "Kragujevac", NameEn: "Kragujevac", NameCyr: "Крагујевац",
	},
	{
		Id: 64, ZpcgStopId: 59, Type: 1,
		Name: "Lapovo", NameEn: "Lapovo", NameCyr: "Лапово",
	},
	{
		Id: 65, ZpcgStopId: 60, Type: 1,
		Name: "Velika Plana", NameEn: "Velika Plana", NameCyr: "Велика Плана",
	},
	{
		Id: 66, ZpcgStopId: 63, Type: 1,
		Name: "Branešci", NameEn: "Branesci", NameCyr: "Бранешци",
	},
	{
		Id: 67, ZpcgStopId: 64, Type: 1,
		Name: "Brodarevo", NameEn: "Brodarevo", NameCyr: "Бродарево",
	},
	{
		Id: 68, ZpcgStopId: 67, Type: 1,
		Name: "Vrbnica", NameEn: "Vrbnica", NameCyr: "Врбница",
	},
	{
		Id: 69, ZpcgStopId: 68, Type: 1,
		Name: "Zmajevo", NameEn: "Zmajevo", NameCyr: "Змајево",
	},
	{
		Id: 70, ZpcgStopId: 69, Type: 1,
		Name: "Inđija", NameEn: "Indjija", NameCyr: "Инђија",
	},
	{
		Id: 71, ZpcgStopId: 70, Type: 1,
		Name: "Stara Pazova", NameEn: "Stara Pazova", NameCyr: "Стара Пазова",
	},
	{
		Id: 72, ZpcgStopId: 71, Type: 1,
		Name: "Nova Pazova", NameEn: "Nova Pazova", NameCyr: "Нова Пазова",
	},
	{
		Id: 73, ZpcgStopId: 72, Type: 1,
		Name: "Lovćenac", NameEn: "Lovcenac", NameCyr: "Ловћенац",
	},
	{
		Id: 74, ZpcgStopId: 73, Type: 2,
		Name: "Beška", NameEn: "Beska", NameCyr: "Бешка",
	},
	{
		Id: 75, ZpcgStopId: 74, Type: 2,
		Name: "Bačka Topola", NameEn: "Backa Topola", NameCyr: "Бачка Топола",
	},

	// These are cities near the railway network that have no active train service.
	// They are present so tests can verify the matcher does NOT surface them
	// when the official station index is searched with a strict threshold.
	{
		Id: -1, Blacklisted: true,
		Name: "Budva", NameCyr: "Будва",
	},
	{
		Id: -2, Blacklisted: true,
		Name: "Tivat", NameCyr: "Тиват",
	},
	{
		Id: -3, Blacklisted: true,
		Name: "Kotor", NameCyr: "Котор",
	},
	{
		Id: -4, Blacklisted: true,
		Name: "Cetinje", NameCyr: "Цетине",
	},
	{
		Id: -5, Blacklisted: true,
		Name: "Perast", NameCyr: "Пераст",
	},
	{
		Id: -6, Blacklisted: true,
		Name: "Durmitor", NameCyr: "Дурмитор",
	},
	{
		Id: -7, Blacklisted: true,
		Name: "Petrovac", NameCyr: "Петровац",
	},
	{
		Id: -8, Blacklisted: true,
		Name: "Ulcinj", NameCyr: "Улцињ",
		ProductionAliases: []string{"Ульцин"},
	},
	{
		Id: -9, Blacklisted: true,
		Name: "Sveti Stefan", NameCyr: "Свети Стефан",
	},
	{
		Id: -10, Blacklisted: true,
		Name: "Becici", NameCyr: "Бечичи",
	},
	{
		Id: -11, Blacklisted: true,
		Name: "Herceg Novi", NameCyr: "Херцег Нови",
	},
	{
		Id: -12, Blacklisted: true,
		Name: "Savnik", NameCyr: "Шавник",
	},
	{
		Id: -13, Blacklisted: true,
		Name: "Zabljak", NameCyr: "Жабляк",
	},
	{
		Id: -14, Blacklisted: true,
		Name: "Albania", NameCyr: "Албания",
	},
	{
		Id: -15, Blacklisted: true,
		Name: "Tirana", NameCyr: "Тирана",
	},
	{
		Id: -16, Blacklisted: true,
		Name: "Shkoder", NameCyr: "Шкодер",
	},
	{
		Id: -17, Blacklisted: true,
		Name: "Bosnia and Herzegovina", NameCyr: "Босния и Герцеговина",
	},
	{
		Id: -18, Blacklisted: true,
		Name: "Sarajevo", NameCyr: "Сараево",
	},
}
