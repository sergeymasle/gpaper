// Данная библиотека обеспечивает простую транслитерацию. Для этого достаточно только описать
// соответствующую таблицу подстановки символов, в которой буква проассоциирована с какой-либо
// строкой. В частности, такая ассоциация для транслитерации русского языка уже определена.
//
// Текущая реализация подразумевает только одностороннюю транслитерацию: обратное преобразование
// будет уже не так очевидно.
//
// Хоть кода в этой библиотеке и не очень много, но время на нее все-таки было потрачено, т.к.
// раньше я просто не задумывался о некоторых аспектах работы с транслитерацией.
//
// В общем, как и большинство других аналогичных библиотек, она перебирает все символы в строке
// и заменяет их по предложенному ей словарю. Отличие только в том, что, с моей точки зрения, она
// более корректно отрабатывает случаи с чередованием заглавных букв. Например:
//
//     "ЧАЩА" -> "CHASCHA"
//     "ЧаЩа" -> "ChaScha"
//     "Чаща" -> "Chascha"
//     "чаЩА" -> "chaSCHA"
//
// Для транслитерации русских букв в ней уже предусмотрен встроенный словарь. Для других языков
// вы можете задать свой. Все достаточно просто:
//
//  import "github.com/mdigger/translit"
//
//  tests := []string{
//     "Проверочная СТРОКА для транслитерации",
//     "ЧАЩА",
//     "ЧаЩа",
//     "Чаща",
//     "чаЩА",
//  }
//  for _, text := range tests {
//      fmt.Println(translit.Ru(text))
//  }
package main

import (
	"bytes"
	"strings"
	"unicode"
	"golang.org/x/exp/utf8string"
)

// RuTransiltMap описывает замены русских букв на английские при транслитерации. Некоторые буквы
// заменяются ни на одну, а на две или три буквы латинского алфавита. А мягкий знак вообще исчезает.
// Но такова обычная распространенная схема транслитерации.
var RuTransiltMap = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "yo",
	'ж': "zh",
	'з': "z",
	'и': "i",
	'й': "j",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "h",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "sch",
	'ъ': "'",
	'ы': "y",
	'ь': "",
	'э': "e",
	'ю': "ju",
	'я': "ja",
	' ': "-",
}

// Transliterate выполняет транслитерацию в строке по указанной таблице и возвращает новую строку с
// результатом такого преобразования. Все символы, которые не указаны в таблице транслитерации,
// останутся без изменения.
//
// При транслитерировании учитывается, что замена буквы может быть произведена на строку
// произвольной длины и корректно обрабатываются чередования заглавных и строчных букв. В частности,
// производится корректная транслитерация следующих случаев:
//  "ЧАЩА" -> "CHASCHA"
//  "ЧаЩа" -> "ChaScha"
//  "Чаща" -> "Chascha"
//  "чаЩА" -> "chaSCHA"
//
// При желании, вы можете указать любую таблицу в качестве второго параметра при вызове функции,
// по которой и будет выполнено данное преобразование.
func Transliterate(text string, translitMap map[rune]string) string {
	var result bytes.Buffer
	utf8text := utf8string.NewString(text)
	length := utf8text.RuneCount()
	for index := 0; index < length; index++ {
		runeValue := utf8text.At(index)
		switch str, ok := translitMap[unicode.ToLower(runeValue)]; {
		case !ok:
			result.WriteRune(runeValue)
		case str == "":
			continue
		case unicode.IsUpper(runeValue):
			// Если следующий или предыдущий символ тоже заглавная буква, то все буквы строки
			// заглавные. Иначе, заглавная только первая буква.
			if (length > index+1 && unicode.IsUpper(utf8text.At(index+1))) ||
				(index > 0 && unicode.IsUpper(utf8text.At(index-1))) {
				str = strings.ToUpper(str)
			} else {
				str = strings.Title(str)
			}
			fallthrough
		default:
			result.WriteString(str)
		}
	}
	return result.String()
}

func Transliterate2(text *string, translitMap *map[rune]string) {
	var stringBuffer bytes.Buffer
	mapa := *translitMap
	utf8text := utf8string.NewString(*text)

	var r rune
	var str string
	for i := 0; i < utf8text.RuneCount(); i++ {
		r = utf8text.At(i)

		str = mapa[unicode.ToLower(r)]
		stringBuffer.WriteString(str)
	}
	*text = stringBuffer.String()
}

// Ru выполняет транслитерацию строки с учетом словаря для русской транслитерации.
func Ru(text string) string {
	return Transliterate(text, RuTransiltMap)
}
