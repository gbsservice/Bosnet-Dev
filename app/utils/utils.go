package utils

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/dustin/go-humanize"
	"html/template"
	"log"
	"math"
	"api_kino/config/database"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type HtmlPDFOption struct {
	FileName    string
	FilePath    string
	PageSize    string
	Orientation string
	Data        interface{}
}

type HtmlOption struct {
	FileName string
	FilePath string
	Data     interface{}
}

func ToDate(toRound time.Time) time.Time {
	return time.Date(toRound.Year(), toRound.Month(), toRound.Day(), 0, 0, 0, 0, toRound.Location())
}

func Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793
	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)
	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)
	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}
	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}
	return dist
}

func ToTerbilang(num int) string {
	var s string
	satuan := [12]string{"", "Satu", "Dua", "Tiga", "Empat", "Lima", "Enam", "Tujuh", "Delapan", "Sembilan", "Sepuluh", "Sebelas"}
	if num < 12 {
		s = satuan[num]
	} else if num < 20 {
		s = fmt.Sprintf("%s Belas", ToTerbilang(num-10))
	} else if num < 100 {
		s = fmt.Sprintf("%s Puluh %s", ToTerbilang(num/10), ToTerbilang(num%10))
	} else if num < 200 { // ratus
		s = fmt.Sprintf("Seratus %s", ToTerbilang(num-100))
	} else if num < 1000 {
		s = fmt.Sprintf("%s Ratus %s", ToTerbilang(num/100), ToTerbilang(num%100))
	} else if num < 2000 { // ribu
		s = fmt.Sprintf("Seribu %s", ToTerbilang(num-1000))
	} else if num < 1000000 {
		s = fmt.Sprintf("%s Ribu %s", ToTerbilang(num/1000), ToTerbilang(num%1000))
	} else if num < 2000000 { // juta
		s = fmt.Sprintf("Satu Juta %s", ToTerbilang(num-1000000))
	} else if num < 1000000000 {
		s = fmt.Sprintf("%s Juta %s", ToTerbilang(num/1000000), ToTerbilang(num%1000000))
	} else if num < 2000000000 { // milyar
		s = fmt.Sprintf("Satu Milyar %s", ToTerbilang(num-1000000000))
	} else if num < 1000000000000 {
		s = fmt.Sprintf("%s Milyar %s", ToTerbilang(num/1000000000), ToTerbilang(num%1000000000))
	} else if num < 2000000000000 { // triliun
		s = fmt.Sprintf("Satu Triliun %s", ToTerbilang(num-1000000000000))
	} else if num < 1000000000000000 {
		s = fmt.Sprintf("%s Triliun %s", ToTerbilang(num/1000000000000), ToTerbilang(num%1000000000000))
	}
	return strings.TrimSpace(s)
}

func toTerbilangRp(num float64) string {
	return fmt.Sprintf("%s rupiah", ToTerbilang(int(num)))
}

func toTerbilangSuffix(suff string, num float64) string {
	return fmt.Sprintf("%s %s", ToTerbilang(int(num)), suff)
}

func getProp(d interface{}, label string) (interface{}, bool) {
	switch reflect.TypeOf(d).Kind() {
	case reflect.Struct:
		v := reflect.ValueOf(d).FieldByName(label)
		return v.Interface(), true
	}
	return nil, false
}

func lang(in string) string {
	r := strings.NewReplacer(
		"January", "Januar",
		"February", "Februar",
		"March", "März",
		"April", "April",
		"May", "Mai",
		"June", "Juni",
		"July", "Juli",
		"August", "August",
		"September", "September",
		"October", "Oktober",
		"November", "November",
		"December", "Dezember",
	)
	return r.Replace(in)
}

func sum(t interface{}, fieldName string) float64 {
	val := From(t).
		SelectT(func(u interface{}) float64 {
			v := reflect.ValueOf(u).FieldByName(fieldName)
			return v.Float()
		}).
		SumFloats()
	return val
}

func isNull(val *string, replace string) string {
	if val == nil {
		return replace
	}
	return *val
}

func add(x, y int) int {
	return x + y
}

func HtmlPDF(option HtmlPDFOption) ([]byte, error) {
	path, err := os.Getwd()
	t, err := template.New(option.FileName).Funcs(template.FuncMap{
		"currency":          humanize.Commaf,
		"isNull":            isNull,
		"toTerbilangRp":     ToTerbilang,
		"toTerbilangSuffix": toTerbilangSuffix,
		"add":               add,
		"sum":               sum,
		"lang":              lang,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	}).ParseFiles(filepath.Join(path, option.FilePath))
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, option.Data); err != nil {
		return nil, err
	}
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(tpl.Bytes()))
	pdf.AddPage(page)
	if option.PageSize != "" {
		pdf.PageSize.Set(option.PageSize)
	} else {
		pdf.PageSize.Set(wkhtmltopdf.PageSizeA5)
	}
	if option.Orientation != "" {
		pdf.Orientation.Set(option.Orientation)
	} else {
		pdf.Orientation.Set(wkhtmltopdf.OrientationLandscape)
	}
	pdf.Dpi.Set(300)

	//if option.Writer != nil {
	//	pdf.SetOutput(option.Writer)
	//}
	if err := pdf.Create(); err != nil {
		return nil, err
	}
	return pdf.Bytes(), err
}

func Html(option HtmlOption) ([]byte, error) {
	path, err := os.Getwd()
	t, err := template.New(option.FileName).Funcs(template.FuncMap{
		"currency":          humanize.Commaf,
		"isNull":            isNull,
		"toTerbilangRp":     ToTerbilang,
		"toTerbilangSuffix": toTerbilangSuffix,
		"add":               add,
		"sum":               sum,
		"lang":              lang,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	}).ParseFiles(filepath.Join(path, option.FilePath))
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, option.Data); err != nil {
		return nil, err
	}
	return tpl.Bytes(), nil
}

func QueryRun(query string, values ...interface{}) ([]map[string]interface{}, error) {
	db := database.DB
	q := db
	s := strings.Split(query, ",")
	for i := range s {
		log.Println(s[i])
	}

	log.Println(s)
	rows, err := q.Raw(query, values...).Rows()
	if err != nil {
		return nil, err
	}
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	colNum := len(columns)
	var results []map[string]interface{}
	for rows.Next() {
		r := make([]interface{}, colNum)
		for i := range r {
			r[i] = &r[i]
		}
		err = rows.Scan(r...)
		if err != nil {
			return nil, err
		}
		var row = map[string]interface{}{}
		for i := range r {
			row[columns[i]] = r[i]
		}
		results = append(results, row)
	}
	return results, nil
}
