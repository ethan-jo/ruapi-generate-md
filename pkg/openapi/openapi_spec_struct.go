package openapi

type OpenApiSpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    OpenApiInfo            `json:"info"`
	Tags    []OpenApiTag           `json:"tags"`
	Paths   map[string]OpenApiPath `json:"paths"`
}

type OpenApiInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type OpenApiTag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type OpenApiPath struct {
	Summary     string            `json:"summary,omitempty"`
	Description string            `json:"description,omitempty"`
	Get         *OpenApiOperation `json:"get,omitempty"`
	Post        *OpenApiOperation `json:"post,omitempty"`
	Put         *OpenApiOperation `json:"put,omitempty"`
	Delete      *OpenApiOperation `json:"delete,omitempty"`
	Patch       *OpenApiOperation `json:"patch,omitempty"`
	Options     *OpenApiOperation `json:"options,omitempty"`
	Head        *OpenApiOperation `json:"head,omitempty"`
	Trace       *OpenApiOperation `json:"trace,omitempty"`
}

type OpenApiOperation struct {
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Parameters  []OpenApiParameter         `json:"parameters,omitempty"`
	RequestBody OpenApiRequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]OpenApiResponse `json:"responses,omitempty"`
}

type OpenApiParameter struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	In          string      `json:"in,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
}

type OpenApiResponse struct {
	Description string                      `json:"description,omitempty"`
	Content     map[string]OpenApiMediaType `json:"content,omitempty"`
}

type OpenApiRequestBody struct {
	Content map[string]OpenApiMediaType `json:"content,omitempty"`
}

type OpenApiMediaType struct {
	Schema  OpenApiSchema `json:"schema,omitempty"`
	Example string        `json:"example,omitempty"`
}

type OpenApiSchema struct {
	Type       string                           `json:"type,omitempty"`
	Properties map[string]OpenApiSchemaProperty `json:"properties,omitempty"`
}

type OpenApiSchemaProperty struct {
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}
