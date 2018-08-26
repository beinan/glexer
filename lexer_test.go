package glexer

import "testing"

var simpleSchema = `
  # A simple type
  # for lexing test
  type AType{
    id ID
    list: [Int!]
  }
`
var schema = `
  enum Gender {
    MALE
    FEMALE
  }
  interface IDNode{
    id: ID!
  }
  type Edge{
    node: IDNode!
    cursor: String!
  }
  type Edges {
    edges: [Edge]
    hasMore: Boolean!
  }
  input EdgesInput {
    cursor: String
    pageSize: Int!
    isRev: Boolean #is reversed ordering
  }
  type User implements IDNode{
    id: ID!
    name: String!
    gender: Gender
    friends(pageNum: Int = 0, pageSize: Int = 20): [User!]
    friendEdges(input: EdgesInput): Edges
  }
  input AuthInput {
    loginName: String!
    password: String!
  }
    
  type Query {
    getUser(id: ID!): User  
  }
  type Mutation {
    signUp(input: AuthInput!): User
    signIn(input: AuthInput!): String!
    addFriend(fromId: ID!, toId: ID!): Boolean!
  }
  schema {
		query: Query
		mutation: Mutation
	}
`

func TestParseSchema(t *testing.T) {
	err := ParseSchema(simpleSchema)
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}

var query = `
{
  getUser(id: "5d708428010004") {
    id
    name
    friendEdges(input: {cursor: "aa", pageSize: 5}) {
      edges {
        node {
          ... on User {
            name
          }
        }
      }
      hasMore
    }
  }
}
`

func TestParseQuery(t *testing.T) {
	err := ParseQuery(query)
	if err != nil {
		t.Errorf("Parse query failed: %v", err)
	}
}
