package spellcheck

import (
	"net/url"
	"testing"
)

func TestCheck(t *testing.T) {
	checker := NewYandexSpeller()
	result, err := checker.Check("Привет мир! Hello world!")
	if len(result) != 0 || err != nil {
		t.Fail()
	}

	bigText := `Larian Studios значительно увеличила штат сотрудников для работы над 
        игрой — к весне 2020 года в студии работало 250 человек, и созданием игры также 
        занималось ещё 50 внешних сотрудников, привлечённых в рамках аутсорсинга[15]. 
        К моменту выхода игры в 2023 году штат студии вырос до 450 человек[16]. При этом 
        компания, будучи и разработчиком, и издателем игры, сохраняет творческую независимость 
        и не обязана отчитываться ни перед кем, кроме Wizards of the Coast[15]. 
        
        Baldur’s Gate 3 
        должна быть более мрачной и жестокой, чем и предыдущие игры Baldur’s Gate, и серия Divinity; 
        разработчики из Larian считали, что рискуют, предлагая Wizards of the Coast придуманную ими 
        концепцию для игры, однако правообладатель посчитал предложение «крутым»[7].
        `
	result, err = checker.Check(bigText)
	if len(result) != 0 || err != nil {
		t.Fail()
	}

	urlEncoded := url.QueryEscape(bigText)
	result, err = checker.Check(urlEncoded)
	if len(result) != 0 || err != nil {
		t.Fail()
	}
}

func TestCheckBad(t *testing.T) {
	checker := NewYandexSpeller()
	result, err := checker.Check("Привет мирп! Hello wrld!")
	if len(result) != 2 || err != nil {
		t.Fail()
	}
	text := "Привет мирп! Hello wrld!"
	urlEncoded := url.QueryEscape(text)
	result, err = checker.Check(urlEncoded)
	if len(result) != 2 || err != nil {
		t.Fail()
	}

	bigText2Mistakes := `Larian Studios значително << HERE << увеличила штат сотрудникофф << HERE << 
        для работы над игрой — к весне 2020 года в студии работало 250 человек, и созданием игры 
        также занималось ещё 50 внешних сотрудников, привлечённых в рамках аутсорсинга[15]. 
        К моменту выхода игры в 2023 году штат студии вырос до 450 человек[16]. При этом 
        компания, будучи и разработчиком, и издателем игры, сохраняет творческую независимость 
        и не обязана отчитываться ни перед кем, кроме Wizards of the Coast[15]. 
        
        Baldur’s Gate 3 
        должна быть более мрачной и жестокой, чем и предыдущие игры Baldur’s Gate, и серия Divinity; 
        разработчики из Larian считали, что рискуют, предлагая Wizards of the Coast придуманную ими 
        концепцию для игры, однако правообладатель посчитал предложение «крутым»[7].
        `
	result, err = checker.Check(bigText2Mistakes)
	if len(result) != 2 || err != nil {
		t.Fail()
	}

	urlEncoded = url.QueryEscape(bigText2Mistakes)
	result, err = checker.Check(urlEncoded)
	if len(result) != 2 || err != nil {
		t.Fail()
	}
}
