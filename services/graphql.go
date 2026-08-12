package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GraphQLService handles GraphQL introspection and query execution.
type GraphQLService struct{}

type GraphQLRequest struct {
	URL       string            `json:"url"`
	Query     string            `json:"query"`
	Variables string            `json:"variables"`
	Headers   []Pair            `json:"headers"`
}

type GraphQLResponse struct {
	Data     string `json:"data"`
	Errors   string `json:"errors,omitempty"`
	Duration int64  `json:"duration"`
	Error    string `json:"error,omitempty"`
}

type GraphQLSchema struct {
	Types    []GraphQLType    `json:"types"`
	Queries  []GraphQLField   `json:"queries"`
	Mutations []GraphQLField  `json:"mutations"`
}

type GraphQLType struct {
	Name   string         `json:"name"`
	Kind   string         `json:"kind"`
	Fields []GraphQLField `json:"fields,omitempty"`
}

type GraphQLField struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Type        string              `json:"type"`
	Args        []GraphQLArgument   `json:"args,omitempty"`
}

type GraphQLArgument struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

// ExecuteQuery sends a GraphQL query/mutation.
func (g *GraphQLService) ExecuteQuery(req GraphQLRequest) GraphQLResponse {
	start := time.Now()

	// Build request body
	variables := map[string]interface{}{}
	if req.Variables != "" {
		json.Unmarshal([]byte(req.Variables), &variables)
	}

	payload := map[string]interface{}{
		"query":     req.Query,
		"variables": variables,
	}
	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequest("POST", req.URL, bytes.NewReader(body))
	if err != nil {
		return GraphQLResponse{Error: fmt.Sprintf("create request: %v", err), Duration: time.Since(start).Milliseconds()}
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for _, h := range req.Headers {
		if h.Enabled && h.Name != "" {
			httpReq.Header.Set(h.Name, h.Value)
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return GraphQLResponse{Error: fmt.Sprintf("request failed: %v", err), Duration: time.Since(start).Milliseconds()}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	duration := time.Since(start).Milliseconds()

	// Parse GraphQL response
	var gqlResp struct {
		Data   json.RawMessage `json:"data"`
		Errors json.RawMessage `json:"errors"`
	}
	json.Unmarshal(respBody, &gqlResp)

	result := GraphQLResponse{Duration: duration}
	if gqlResp.Data != nil {
		formatted, _ := json.MarshalIndent(gqlResp.Data, "", "  ")
		result.Data = string(formatted)
	}
	if gqlResp.Errors != nil {
		formatted, _ := json.MarshalIndent(gqlResp.Errors, "", "  ")
		result.Errors = string(formatted)
	}
	if result.Data == "" && result.Errors == "" {
		// Couldn't parse as GraphQL response; return raw
		result.Data = string(respBody)
	}

	return result
}

// IntrospectSchema fetches the GraphQL schema via introspection query.
func (g *GraphQLService) IntrospectSchema(url string, headers []Pair) (GraphQLSchema, error) {
	introspectionQuery := `{
		__schema {
			types {
				name
				kind
				fields {
					name
					description
					type { name kind ofType { name kind ofType { name kind } } }
					args { name type { name kind ofType { name kind } } defaultValue }
				}
			}
			queryType { name }
			mutationType { name }
		}
	}`

	req := GraphQLRequest{
		URL:     url,
		Query:   introspectionQuery,
		Headers: headers,
	}

	resp := g.ExecuteQuery(req)
	if resp.Error != "" {
		return GraphQLSchema{}, fmt.Errorf("%s", resp.Error)
	}

	// Parse introspection result
	var introspection struct {
		Schema struct {
			Types []struct {
				Name   string `json:"name"`
				Kind   string `json:"kind"`
				Fields []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Type        struct {
						Name   *string `json:"name"`
						Kind   string  `json:"kind"`
						OfType *struct {
							Name *string `json:"name"`
							Kind string  `json:"kind"`
						} `json:"ofType"`
					} `json:"type"`
					Args []struct {
						Name string `json:"name"`
						Type struct {
							Name   *string `json:"name"`
							Kind   string  `json:"kind"`
							OfType *struct {
								Name *string `json:"name"`
								Kind string  `json:"kind"`
							} `json:"ofType"`
						} `json:"type"`
						DefaultValue *string `json:"defaultValue"`
					} `json:"args"`
				} `json:"fields"`
			} `json:"types"`
			QueryType    *struct{ Name string } `json:"queryType"`
			MutationType *struct{ Name string } `json:"mutationType"`
		} `json:"__schema"`
	}

	// The data from ExecuteQuery is already the "data" field
	json.Unmarshal([]byte(resp.Data), &introspection)

	schema := GraphQLSchema{}

	for _, t := range introspection.Schema.Types {
		if len(t.Name) > 0 && t.Name[0] == '_' {
			continue // Skip internal types
		}

		gqlType := GraphQLType{Name: t.Name, Kind: t.Kind}
		for _, f := range t.Fields {
			typeName := resolveTypeName(f.Type.Name, f.Type.Kind, f.Type.OfType)
			field := GraphQLField{
				Name:        f.Name,
				Description: f.Description,
				Type:        typeName,
			}
			for _, a := range f.Args {
				argType := resolveTypeName(a.Type.Name, a.Type.Kind, a.Type.OfType)
				dv := ""
				if a.DefaultValue != nil {
					dv = *a.DefaultValue
				}
				field.Args = append(field.Args, GraphQLArgument{
					Name:         a.Name,
					Type:         argType,
					DefaultValue: dv,
				})
			}
			gqlType.Fields = append(gqlType.Fields, field)
		}
		schema.Types = append(schema.Types, gqlType)

		// Identify query/mutation fields
		if introspection.Schema.QueryType != nil && t.Name == introspection.Schema.QueryType.Name {
			schema.Queries = gqlType.Fields
		}
		if introspection.Schema.MutationType != nil && t.Name == introspection.Schema.MutationType.Name {
			schema.Mutations = gqlType.Fields
		}
	}

	return schema, nil
}

func resolveTypeName(name *string, kind string, ofType *struct {
	Name *string `json:"name"`
	Kind string  `json:"kind"`
}) string {
	if name != nil {
		return *name
	}
	if ofType != nil && ofType.Name != nil {
		switch kind {
		case "NON_NULL":
			return *ofType.Name + "!"
		case "LIST":
			return "[" + *ofType.Name + "]"
		default:
			return *ofType.Name
		}
	}
	return kind
}
