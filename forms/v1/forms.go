package forms

import (
	forms "google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

func ScopesAll() []string {
	return []string{
		forms.DriveScope,
		forms.DriveFileScope,
		forms.FormsBodyScope,
		forms.FormsResponsesReadonlyScope}
}

func ClientOptionScopesAll() option.ClientOption {
	scopes := ScopesAll()
	return option.WithScopes(scopes...)
}

func ScopesReadOnly() []string {
	return []string{
		forms.DriveReadonlyScope,
		forms.FormsBodyReadonlyScope,
		forms.FormsResponsesReadonlyScope}
}

func ClientOptionScopesReadOnly() option.ClientOption {
	scopes := ScopesReadOnly()
	return option.WithScopes(scopes...)
}
