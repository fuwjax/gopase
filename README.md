## Go Parsing and Serialization

Does the world need another parser? No, of course not. So, why another parser?

Two reasons: I wanted a fun way to learn Go, and I want to have something I can point to while I explain why parsing makes every developer's life better.

### The Joys of Parsing

Parsing gets a bad rap. Parser generators are confusing, often requiring an archaic form of mental gymnastics using obscure languages to produce grammars which are difficult to write, challenging to read, and impossible to maintain. So we stick to the parsers someone else writes, but far too often those are tied into a whole host of things we don't necessarily want, like object mapping and traversal, resource allocation, and threading limitations.

Why does it need to be this way?

Parsing is just a function that takes a byte or character stream and transforms it into something else. That something else can be pretty broad and still qualify as a parser. It can be a specific object, like the DOM from an HTML parser, or more generic like arrays/maps from a generic JSON parser. It can even create components in your application or do nothing more than trigger callbacks to some handler in the case of SAX. Most of the time it's a rather useless, temporal structure called an Abstract Syntax Tree (AST). So a parser just takes data and turns it into something, that's a pretty broad classification for a function.

There are loads of algorithms available to the ambitious parser writer. For this project, I intend to implement 3 - PEG, Earley and TBD. That last one isn't a parser algorithm, it's just me being wishy-washy in what sort of top-down parser I'll add to the mix. Gonna be honest here, I know they're the only actual popular parser strategies, and I understand the reasons why. But if you are writing your own parser without wanting to become a parser fanatic, then you almost certainly want to use PEG.

### So what makes gopase unique

Three things, I think. First, gopase is an OO approach to a traditionally not-OO problem. We'll talk about why that matters at some point, but for now we'll just agree that it's somewhat unique and move on. 

Second, the "something" gopase returns is whatever you want. It uses a ridiculously simple callback strategy in each of the parsers that takes almost all the guesswork out of producing a meaningful object representation of the parse.

Third, and perhaps most useful, gopase is intended to be an illustration of what parsing actually brings to the table - why these different algorithms are useful, and what you can do with them.

### So what is a parser

Parsing, for the purposes of gopase at least, requires 5 things. The most obvious is the input stream, the sequence of characters or bytes that the parser will consume to produce the output. The second is the grammar - the specification defining the rules that input stream must follow. The third is the root, which of those rules should govern the start of the parse. Next is something that will handle turning parse results into something meaningful; gopase refers to this as the handler. Finally is the actual parsing engine; think of this as a black box that takes those other 4 things and produces a something. 

Most of the time, when we talk about a parser we mean something that combines the engine and the grammar. So we treat the engine as a parser that consumes a grammar according to a parser-grammar and returns a new parser for the grammar... that all gets very confusing very quickly.

Regular expressions often have a way of matching that includes the expression, something like Regex.match(/[a-f]/, "hello") and something that lets you precompile the expression, so something like Regex.compile(/[a-f]/).match("hello"). Everyone thinks they need the expression at some point, and no one thinks of the compile as something special. It's for efficiency and performance, not some sort of ivory tower magic.

So, parsing requires those 5 things, perhaps not all at the same time, but if you don't have those 5 things, then you had to do something that took a lot of effort to solve a problem that already took a lot of effort.

The normal process is to hand the engine the set of grammar, root, and handler objects and get back a reusable parser that specifically parses those sematics into those outputs. And then you can hand whatever input you want to that reusable parser and it will happily churn out the outputs according to its particular rules and callbacks. So just like with the regex example, we have a notion of a compile step that takes the reusable bits and squishes them together, returning that function from an input to an output.

### Timeline of gopase

1) Bootstrap PEG parser. This is a hardcoded parser that converts grammars into parsers. At this point we have a grammar specification for PEG grammars, and a handler that produces PEG parsers.
2) Template engine. This is a parser that consumes according to a template grammar. So this is a grammar specification for the templating language, and a handler that produces an object that knows how to walk an object graph and render text.
3) PEG code template. This is a template specific to serializing PEG grammars. So this is a template that walks a PEG parser and serializes that into valid go code. At this point the bootstrap can be replaced with this output as desired. The template engine and various samples like JSON and CSV can get their own hardcoded variants.
4) Earley parser. This is a PEG grammar that can produce an Earley parser. This gets a little mind-bendy, but there isn't anything anywhere that says the same strategy has to be used for every turtle. At this point we can start making comparisons between the two.
5) Chesur. I created a language many moons ago that made it super easy to write parsers/serializers with a PEG mindset. I look forward to revisiting it.
6) TBD parser. My gut says this will be a predictive recursive descent parser. Time will tell.

### How do we use this PEG parser

Hand Peg() a string that follows the EBNF-ish extensions to basic PEG grammar given below, and it'll return you a Grammar that conforms to that specification. Then from that Grammar, call Parser() with a root and a handler and you'll have a Parser of your very own.

    grammar := Peg(`
        S = "a" [b-c]
    `)
    parser := grammar.Parser("S", handler)

The PEG specification works like this

    Basic Expressions
        "some string" - Literal match of the exact string sequence, without the double quotes
        'some string' - Literal match of the exact string sequence, without the single quotes
        [a-b] - Regex match of the character class
        . - Matches any single character
        Name - reference match of rule by name, can cycle or recurse
    Extended Expressions (Expr stands for any valid sub expression)
        x y ... - sequence of expressions must follow in that order
        x / y / ... - options for expressions matches the first success
        (x) - expressions can be enclosed in parentheses for grouping
        x? - zero or one times
        x* - zero or more times
        x+ - one or more times
        &x - zero match positive lookahead
        !x - zero match negative lookahead

The handlers are pretty easy. A handler is a struct with a set of public methods. Each Rule that should be handled gets a method
of the same name. This method takes as an argument an iter.Seq2[string, any], effectively a sequence of key-value pairs where the
keys are the Reference names on the right side of that Rule, along with the objects returned from their handlers.

As an example, say we have a rule

    Record = Field (", " Field)* EOL

And a handler

    func (h *MyHandler) Record(results iter.Seq2[string, any]) (any, error) {
        ...
    }

    func (h *MyHandler) Field(results iter.Seq2[string, any]) (any, error) {
        ...
    }

Then the results passed into Record would be a series of "Field" keys with the corresponding result from the Field() handler followed by a single "EOL" result. The EOL rule would
likely not have a handler, and would therefore contain the string match accepted by the Rule, say the string "\n".

