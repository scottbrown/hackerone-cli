package hackeronecli

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

const (
	FormatJSON     = "json"
	FormatText     = "text"
	FormatMarkdown = "markdown"
)

func ParseFormat(s string) (string, error) {
	switch strings.ToLower(s) {
	case FormatJSON:
		return FormatJSON, nil
	case FormatText:
		return FormatText, nil
	case FormatMarkdown:
		return FormatMarkdown, nil
	default:
		return "", fmt.Errorf("unsupported format %q: must be one of json, text, markdown", s)
	}
}

func FormatOutput(w io.Writer, format string, data interface{}) error {
	switch format {
	case FormatText:
		return formatText(w, data)
	case FormatMarkdown:
		return formatMarkdown(w, data)
	default:
		return formatJSON(w, data)
	}
}

func FormatMessage(w io.Writer, format string, msg string) error {
	switch format {
	case FormatJSON:
		return formatJSON(w, map[string]string{"message": msg})
	default:
		_, err := fmt.Fprintln(w, msg)
		return err
	}
}

func formatJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func formatText(w io.Writer, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			_, err := fmt.Fprintln(w, "<nil>")
			return err
		}
		v = v.Elem()
	}

	switch {
	case v.Kind() == reflect.String:
		_, err := fmt.Fprintln(w, v.String())
		return err
	case v.Kind() == reflect.Slice:
		return formatTextTable(w, v)
	case v.Kind() == reflect.Map:
		return formatTextMap(w, v)
	case v.Kind() == reflect.Struct:
		return formatTextSingle(w, v)
	default:
		return formatJSON(w, data)
	}
}

func formatTextSingle(w io.Writer, v reflect.Value) error {
	pairs := extractFields(v)
	for _, p := range pairs {
		if _, err := fmt.Fprintf(w, "%s: %s\n", p.key, p.value); err != nil {
			return err
		}
	}
	return nil
}

func formatTextTable(w io.Writer, v reflect.Value) error {
	if v.Len() == 0 {
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}
	if first.Kind() == reflect.Interface {
		first = first.Elem()
	}

	if first.Kind() == reflect.Map {
		return formatTextMapSlice(w, v)
	}

	if first.Kind() != reflect.Struct {
		for i := 0; i < v.Len(); i++ {
			if _, err := fmt.Fprintln(w, fmt.Sprint(v.Index(i).Interface())); err != nil {
				return err
			}
		}
		return nil
	}

	headers := extractHeaders(first)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		pairs := extractFields(elem)
		vals := make([]string, len(pairs))
		for j, p := range pairs {
			vals[j] = p.value
		}
		if _, err := fmt.Fprintln(tw, strings.Join(vals, "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func formatTextMapSlice(w io.Writer, v reflect.Value) error {
	allKeys := orderedMapKeys(v)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(allKeys, "\t")); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		vals := make([]string, len(allKeys))
		for j, key := range allKeys {
			mv := elem.MapIndex(reflect.ValueOf(key))
			if mv.IsValid() {
				vals[j] = fmt.Sprint(mv.Interface())
			}
		}
		if _, err := fmt.Fprintln(tw, strings.Join(vals, "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func formatTextMap(w io.Writer, v reflect.Value) error {
	keys := make([]string, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, fmt.Sprint(k.Interface()))
	}
	sort.Strings(keys)
	for _, k := range keys {
		val := v.MapIndex(reflect.ValueOf(k))
		if _, err := fmt.Fprintf(w, "%s: %s\n", k, fmt.Sprint(val.Interface())); err != nil {
			return err
		}
	}
	return nil
}

func formatMarkdown(w io.Writer, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			_, err := fmt.Fprintln(w, "*nil*")
			return err
		}
		v = v.Elem()
	}

	switch {
	case v.Kind() == reflect.String:
		_, err := fmt.Fprintln(w, v.String())
		return err
	case v.Kind() == reflect.Slice:
		return formatMarkdownTable(w, v)
	case v.Kind() == reflect.Map:
		return formatMarkdownMap(w, v)
	case v.Kind() == reflect.Struct:
		return formatMarkdownSingle(w, v)
	default:
		return formatJSON(w, data)
	}
}

func formatMarkdownSingle(w io.Writer, v reflect.Value) error {
	pairs := extractFields(v)
	if _, err := fmt.Fprintln(w, "| Key | Value |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- |"); err != nil {
		return err
	}
	for _, p := range pairs {
		if _, err := fmt.Fprintf(w, "| %s | %s |\n", p.key, p.value); err != nil {
			return err
		}
	}
	return nil
}

func formatMarkdownTable(w io.Writer, v reflect.Value) error {
	if v.Len() == 0 {
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}
	if first.Kind() == reflect.Interface {
		first = first.Elem()
	}

	if first.Kind() == reflect.Map {
		return formatMarkdownMapSlice(w, v)
	}

	if first.Kind() != reflect.Struct {
		for i := 0; i < v.Len(); i++ {
			if _, err := fmt.Fprintln(w, fmt.Sprint(v.Index(i).Interface())); err != nil {
				return err
			}
		}
		return nil
	}

	headers := extractHeaders(first)
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(headers, " | ")); err != nil {
		return err
	}
	seps := make([]string, len(headers))
	for i := range seps {
		seps[i] = "---"
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | ")); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		pairs := extractFields(elem)
		vals := make([]string, len(pairs))
		for j, p := range pairs {
			vals[j] = p.value
		}
		if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(vals, " | ")); err != nil {
			return err
		}
	}
	return nil
}

func formatMarkdownMapSlice(w io.Writer, v reflect.Value) error {
	allKeys := orderedMapKeys(v)
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(allKeys, " | ")); err != nil {
		return err
	}
	seps := make([]string, len(allKeys))
	for i := range seps {
		seps[i] = "---"
	}
	if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | ")); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		vals := make([]string, len(allKeys))
		for j, key := range allKeys {
			mv := elem.MapIndex(reflect.ValueOf(key))
			if mv.IsValid() {
				vals[j] = fmt.Sprint(mv.Interface())
			}
		}
		if _, err := fmt.Fprintf(w, "| %s |\n", strings.Join(vals, " | ")); err != nil {
			return err
		}
	}
	return nil
}

func formatMarkdownMap(w io.Writer, v reflect.Value) error {
	keys := make([]string, 0, v.Len())
	for _, k := range v.MapKeys() {
		keys = append(keys, fmt.Sprint(k.Interface()))
	}
	sort.Strings(keys)
	if _, err := fmt.Fprintln(w, "| Key | Value |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "| --- | --- |"); err != nil {
		return err
	}
	for _, k := range keys {
		val := v.MapIndex(reflect.ValueOf(k))
		if _, err := fmt.Fprintf(w, "| %s | %s |\n", k, fmt.Sprint(val.Interface())); err != nil {
			return err
		}
	}
	return nil
}

type fieldPair struct {
	key   string
	value string
}

func extractFields(v reflect.Value) []fieldPair {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var pairs []fieldPair
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		tag := sf.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name := jsonFieldName(sf)
		fv := v.Field(i)

		if fv.Kind() == reflect.Struct && sf.Name == "Attributes" {
			pairs = append(pairs, extractFields(fv)...)
			continue
		}

		if fv.Kind() == reflect.Map && sf.Name == "Attributes" {
			mapKeys := make([]string, 0, fv.Len())
			for _, k := range fv.MapKeys() {
				mapKeys = append(mapKeys, fmt.Sprint(k.Interface()))
			}
			sort.Strings(mapKeys)
			for _, k := range mapKeys {
				mv := fv.MapIndex(reflect.ValueOf(k))
				pairs = append(pairs, fieldPair{key: k, value: fmt.Sprint(mv.Interface())})
			}
			continue
		}

		pairs = append(pairs, fieldPair{key: name, value: fmt.Sprint(fv.Interface())})
	}
	return pairs
}

func extractHeaders(v reflect.Value) []string {
	pairs := extractFields(v)
	headers := make([]string, len(pairs))
	for i, p := range pairs {
		headers[i] = p.key
	}
	return headers
}

func jsonFieldName(sf reflect.StructField) string {
	tag := sf.Tag.Get("json")
	if tag == "" {
		return sf.Name
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return sf.Name
	}
	return name
}

func orderedMapKeys(sliceVal reflect.Value) []string {
	keySet := make(map[string]struct{})
	for i := 0; i < sliceVal.Len(); i++ {
		elem := sliceVal.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}
		for _, k := range elem.MapKeys() {
			keySet[fmt.Sprint(k.Interface())] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
