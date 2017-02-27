%{
package main

%}

%union{
  token Token
  block Block
  blocks []Block
  inline Inline
  inlines []Inline
  url string
}

%token<token> TEXT
%token UNORDERED_LIST_MARKER CR LBRACKET RBRACKET LT GT
%type<block> block unordered_list_item unordered_list line quotation
%type<blocks> blocks
%type<inline> inline inline_text inline_http
%type<inlines> inlines
%type<url> url quotation_prefix

%%

blocks:
      block
      {
        $$ = []Block{$1}
        yylex.(*Lexer).result = $$
      }
      | block blocks
      {
        $$ = append([]Block{$1}, $2...)
        yylex.(*Lexer).result = $$
      }

block:
        unordered_list
        {
          $$ = $1
        }
        | quotation
        {
          $$ = $1
        }
        | line
        {
          $$ = $1
        }

line:
    empty_line
    {
      $$ = Line{Inlines: []Inline{}}
    }
    | inlines CR
    {
      $$ = Line{Inlines: $1}
    }

empty_line:
          CR

inlines:
       inline
       {
        $$ = []Inline{$1}
       }
       | inline inlines
       {
        $$ = append([]Inline{$1}, $2...)
       }

inline:
      inline_http
      | inline_text

inline_text:
      TEXT
      {
        $$ = InlineText{Literal: $1.literal}
      }

inline_http:
           LBRACKET url RBRACKET
           {
            $$ = InlineHttp{Url: $2}
           }

url: TEXT
   {
    $$ = $1.literal
   }

unordered_list:
              unordered_list_item
              {
                $$ = UnorderedList{Items: []UnorderedListItem{$1.(UnorderedListItem)}}
              }
              | unordered_list_item unordered_list
              {
                items := $2.(UnorderedList).Items
                list := UnorderedList{Items: append([]UnorderedListItem{$1.(UnorderedListItem)}, items...)}
                $$ = list
              }

unordered_list_item:
                   UNORDERED_LIST_MARKER inlines CR
                   {
                    $$ = UnorderedListItem{Inlines: $2}
                   }

quotation:
         quotation_prefix blocks quotation_suffix
         {
          $$ = Quotation{Cite: $1, Content: $2}
         }

quotation_prefix:
                GT url GT CR
                {
                  $$ = $2
                }
                | GT GT CR
                {
                  $$ = ""
                }

quotation_suffix:
                LT LT CR

%%
