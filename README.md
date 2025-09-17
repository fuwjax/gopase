## Go Parsing and Serialization

Does the world need another parser? No, of course not. So, why another parser?

Two reasons: I wanted a fun way to learn Go, and I want to have something I can point to while I explain why parsing makes every developer's life better.

### The Joys of Parsing

Parsing gets a bad rap. Parser generators are confusing, often requiring an archaic form of mental gymnastics using obscure languages to produce grammars which are difficult to write, challenging to read, and impossible to maintain. So we stick to the parsers someone else writes, but far too often those are tied into a whole host of things we don't necessarily want, like object mapping and traversal, resource allocation, and threading limitations.

Why does it need to be this way?

Parsing is just a function that takes a byte stream and transforms it into something else. That something else can be pretty broad and still qualify as a parser. It can be a specific object, like the DOM from an HTML parser, or more generic like arrays/maps from a generic JSON parser. It can even create components in your application or do nothing more than trigger callbacks to some handler in the case of SAX. Most of the time it's a rather useless, temporal structure called an Abstract Syntax Tree (AST). So a parser just takes a sequence and turns it into a something, that's a pretty broad classification for a function.

It's that breadth that is ultimately the scariest part of parsing. If it's just a black box that turns persistent data into usable data, then what exactly makes that black box do the right thing. For that matter, how do you even know if the black box is doing the right thing? Where do you even start?

There are loads of algorithms available to the ambitious parser writer. For this project, I intend to implement 3 - PEG, Earley and TBD. That last one isn't a parser algorithm, it's just me being wishy-washy in what sort of top-down parser I'll add to the mix. Gonna be honest here, I know they're the only actual popular parser strategies, and I understand the reasons why. But if you are writing your own parser without wanting to become a parser fanatic, then you almost certainly want to use PEG. Well, except when you should use Earley. 

### So what makes gopase unique

Three things, I think. First, gopase is an OO approach to a traditionally not-OO problem. We'll talk about why that matters at some point, but for now we'll just agree that it's somewhat unique and move on. 

Second, the "something" gopase returns is whatever you want. It uses a ridiculously simple callback strategy in each of the parsers that takes almost all the guesswork out of producing a meaningful object representation of the parse.

Third, and perhaps most useful, gopase is intended to be an illustration of what parsing actually brings to the table - why these different algorithms are useful, and what you can do with them.

### So what is a parser

Parsing, for the purposes of gopase at least, requires 5 things. The most obvious is the input stream, the sequence of characters or bytes that the parser will consume to produce the output. The second is the grammar - the specification defining the rules that input stream must follow. The third is the root, which of those rules should govern the start of the parse. Next is something that will handle turning parse results into something meaningful; gopase refers to this as the handler. Finally is the actual parsing engine; think of this as a black box that takes those other 4 things and produces a something. 

Most of the time, when we talk about a parser we mean something that combines the engine and the grammar. So we treat the engine as a parser that consumes a grammar according to a parser-grammar and returns a new parser for the grammar... that all gets very confusing very quickly.

Regular expressions often have a way of matching that includes the expression, something like Regex.match(/[a-f]/, "hello") and something that lets you precompile the expression, so something like Regex.compile(/[a-f]/).match("hello"). Everyone thinks they need the expression at some point, and no one thinks of the compile as something special. It's for efficiency and performance, not some sort of ivory tower magic.

So, parsing requires those 5 things, perhaps not all at the same time, but if you don't have those 5 things, then you had to do something that took a lot of effort to solve a problem that already took a lot of effort. Mustache, I'm looking at you.

The normal process is to hand the engine the set of grammar, root, and handler objects and get back a reusable parser that specifically parses inputs matching the rules of the grammar into outputs produced by the handler. And then you can hand whatever input you want to that reusable parser and it will happily churn out the outputs according to its particular rules and callbacks. So just like with the regex example, we have a notion of a compile step that takes the reusable bits and squishes them together, returning that function from an input to an output.

### The gopase roadmap

1) Bootstrap PEG parser. This is a hardcoded parser that converts grammars into parsers. At this point we have a grammar specification for PEG grammars, and a handler that produces PEG parsers.
1.1) Sample PEG grammars, say the PEG grammar for the PEG flavor matched by the Bootstrap, CSV, JSON, that sort of thing.
2) Template engine. This is a parser that consumes according to a template grammar. So this is a grammar specification for the templating language, and a handler that produces an object that knows how to walk an object graph and render text. My thought was to just implement Mustache. But it wound up being more fun to create something better.
3) PEG code template. This is a template that takes the object oriented output of the Bootstrap and converts it to an imperative parser. This gets a little overwhelming, but at this point the Bootstrap is a code-generated imperative parser that can produce an OO parser for a given grammar input. This OO parser can then be handed off to the template to produce a code-generated imperative parser for that grammar. So we can give the Bootstrap a grammar file that happens to match it's own ruleset, then hand that off to the template engine to generate the code that is in fact, the Bootstrap itself. Egg, meet thy Chicken. The template engine and various samples like JSON and CSV can get their own hardcoded variants. In theory the template engine could also get an imperative corrolary, but I've honestly never thought about solving that before.
4) Earley parser. This is a PEG grammar that can produce an Earley parser. This gets a little mind-bendy, but there isn't anything anywhere that says the same strategy has to be used for every turtle. At this point we can start making comparisons between the two, and diving into their various pros and cons.
5) Chesur. I created a language many moons ago that made it super easy to write parsers/serializers with a PEG mindset without the silliness of EBNF. I look forward to revisiting it.
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
        [a-b] - Regex match of the character class. This is handed directly to the native regexp package, so if Go supports it, so does this parser
        . - Matches any single character
        Name - reference match of rule by name, can cycle or recurse
    Extended Expressions (x and y stand for any valid sub expression)
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

### What was this about a template engine?

I genuinely tried to write Mustache. I got pretty far down the implementation, but ran into a couple snags. The first is that the only way I could think to implement parts of the Mustache grammar was to either post-process or implement look-behinds. Gopase's PEG implementation already has look-aheads, but look-behinds would break the ways I'm able to make the parser memory efficient. I could work around that, I'm pretty sure, but I was trying to write a template engine, not rethink the parser.

The second was a little more difficult. Technically, templating engines that let you change the open/close delimiters from inside the parse aren't really the sort of thing you can do in most rule-based parsing algorithms. One of the cool things about an OO PEG implementation is that, at least in theory, dynamic changes at runtime to the rules is possible. But I can't really think of a good way to do it that would work well every time, and in a way that end users can just change willy-nilly. Templating engines are often used within builds to render things with very little oversight. I'd hate to implement something that could make things go wildly wrong with no way of telling anyone they did.

I'm decently sure that most Mustache implementations are using a whole bunch of regular expression ReplaceAll's under the covers, but that's just a guess.

Anyway, as I was considering how to implement these two features in the parser in ways that didn't offend my delicate sensibilities, I had another thought. My whole claim is that having your own parser means making a thing you want the way you want it. So...

### Happy
#### Templating for the Soul

Does the world need another templating engine? Nope. Necessity may be the mother of invention, but the grandmama is often the raw unadulterated joy of introducing the world to something it hasn't seen before.

	(^ ^) are the open/close delimeters from all (well, most) of the tags in Happy
	(^name^) usually you'll put something in between the carrots to indicate what you'd like to render in place of the tag
	( ^name^ ) whitespace around tags is usually preserved, but a single space between the carrot and the tag 
			will consume all the whitespace in that direction
	(^# comments can contain anything but close tags ^) they can even be multiline
	(^*name^)...(^/^) will render the bit between the start and end tags. 
			We'll need to talk about contexts and keys a bit more for this to make sense.
	(^!name^)...(^/^) will render the bit between the start and end tags, but only if name doesn't exist in the context
	(^="Text")...(^/^) will define an inline partial. We'll talk more partials later
	(^>"Text") will render and insert a partial

And that's it. Well, for right now. There's one more tag type that's a little bit crazy, too crazy to bother implementing at the moment. But that's it - Comments, Values, Sections, Inversions, Partials, and Includes. We'll have to spend a little bit talking about contexts and keys (names and Texts in the summary above). But yeah, pretty simple stuff.

"Isn't that just mustache with kitty emojis instead of handlebars?" Kinda, yeah. I took a lot of inspiration from mustache, since you know, I had implemented a big chunk of it. But there are a few things here that really stand out to me. And one of the big things is being very direct about how Happy interprets the Mustache idea of "logic-less".

There's logic here. It's a computer program after all, that takes computer-y things and turns them into string-y things. And the logic Happy needs to work its magic has to come from the person writing the template, otherwise it would just be serializing data like JSON does. But the "missing" logic is any meaningful way of influencing that magical work via code written in the template itself.

Each of the tags (aside from comments) takes a "key", and that's it. The key is just a way of selecting data from the context. The context is a stack of data, and the only time anything is pushed on or popped off the stack is when entering or leaving a Section (the # tag). Well, key resolution can temporarily use the stack, but don't worry abou tthat just yet.

So a template is just a set of 5 instructions where the only information needed to perform the operation is a key and optionally, some content between a start and end tag. A Value tag (^^) resolves the key against the context to get some value, and then replaces the tag in the output with that value converted to a string. The Section tag (^*^) resolves the key against the context to get some value, and then iterates over the contents of that value by pushing them on to the context one at a time and rendering the content until it reaches the corresponding end tag (^/^). The Inversion tag (^!^) resolves the key against the context to get some value, and if that value is false-y (in go-ese, a zero value) then it renders the bit until the end tag (^/^). The Partial tag (^=^) resolves the key against the context to get some value, turns that value into a string as though it were replacing a Value tag, and then snips the content until the end tag (^/^) and registers it as a partial. The Include tag (^>^) resovles the key, turns it into a string, looks for a partial registered by that name, and replaces the tag with the rendered partial.

Just to say it somewhere, partials don't have to be defined in the template where they're used. They can be passed in as a map to the Render function.

The context is just a stack. The partials are just a map. There are only 5 instructions. How does this even do anything useful?

### Keys are weird, aren't they...

Yup. Keys are a way of selecting data from the stack. Well, they don't have to be, I guess. 

	(^"literal"^) would literally replace the tag with the literal "literal"

Usually, you'd use a naked literal for Partials and Includes, but I'm not here to tell you what to do. Well, I mean, I guess I am. But just this once, you do you.

	(^name^) means search the context for the first element that has a "name" and replace this tag with that "name" value

That element could be a map[string]Something or yourpackage.Person, pointer or raw.

	(^0^) you could even use a number to pull from a slice or array.

You can't get any properties from a scalar (things like numbers or booleans or strings), but the property you pull can be a scalar. Doesn't have to be. Things can get pretty weird.

And it doesn't have to be from the top of the stack. So, if you think about nested Sections, it's going to travel back up through the Section data looking for a match to the key.

Cool. That doesn't feel so weird. 

	(^.^) a dot means to use the top of the stack in its entirety. 
	(^@^) an at means to use the index of the top of the stack from its Section container

Ok, that might be confusing. Remember that a Section is sort of the control flow operation. It tries to resolve the key, and if that key is a container (map, slice, array) then it interates over that container and renders each key-value pair with the Section's inner template. The . means to use the value, the @ means to use the container key. The container key is never included in the search path for key resolution. Which is confusing now that I say it out loud: we don't use keys for figureing out keys. But it's the truth. The only way to get the container key is with the @ and even then, it's only the top level one.

Just to say it, if the section's key resolves to something that isn't a container, it just renders the inner template once with that value. So it acts more like an "if" than a "for" in that situation.

Sometimes you want to go walking around a data structure, and it would be exhausting to have a whole bunch of nested sections. So instead you can just use dots to separate the nested keys.

	(^you.can.even.321.use.numbers^)

So it'll try to resolve "you" against the context, and if it does, it will push it on the stack (temporarily, of course, just for this one key resolution) and then resolve"can" against that new context. If "even" here were to return an array or slice and that array or slice had at least 322 elements, then the 321 would try to grab the 321th element, you know, from zero.

If the key resolves all the way down the dotted sequence, then whatever that final result is will be used to replace the tag, and all the temporary things stuffed into the context will be cleansed from their unholy union... I mean, popped from the top. Key resolution doesn't change the stack for anything other than key resolution.

I know what you're thinking. "You said this was wierd, but really it's just a very long-winded way of saying 'Like Mustache'". And you're right. Did I mention I implemented a mostly mustache template engine before I decided to do something "better"? Better here means by my metric which is explicitly "not having to make a whole bunch of changes to my sacred parser."

Ok, so here's where things get a little bonkers. 

	(^you.can["use"][brackets].too^)

I know, I know... at this point you're wondering if there's someone you should call. Like a wellness check. Do they do wellness checks for inflated egos? Is it ok to ask the cops to swing by and check on someone because they hung their own macaroni art on their fridge? Stick with me, I promise this is a bit brain melting. Yeah, brackets, head exploding brackets.

remember, keys are all about context searching, and key resolution temporarily sticks things on the context stack, yada yada. Brackets don't go on the stack until they're done. So, "you" resolves, goes on the stack. "can["use"][brackets]" resolves and goes on the stack. Wait, what does that mean, exactly?

So, "can" searches the context looking for a property called "can". Finds it, but doesn't push it on to the stack. Then we resolve '"use"' against the stack. The literal string "use". Remember way back at the beginning when I mocked you a bit for using a string literal for something other than partial names? Well, here, we resolve the literal "use" to literally itself... ok, not impressed, I see that. But look at "brackets". We resolve "brackets" against the same context we resolved "can" against. Weird, right?

Because then, after we resolve the value for "can", we take "use" and resolve it against "can" and we resolve the result of resolving brackets against the context against the result of use against can... and I lost you.

Told you this was weird. There are two resolutions, once trying to figure out what "brackets" means against the context of "you" and once more trying to figure out whatever that value was against the context of "you.can.use".

But there's more. Bracket expressions, in addition to being able to work like dotted expressions against maps and structs and interfaces and slices and arrays. Oh... did I mention that you can hit interfaces? Methods on structs? Yup, no problem. You can stick funcs in those maps and slices and arrays too. But anyway, in addition to just fetching things like dotted keys, brackets let you call those functions and methods with arguments. You can call zero arg functions with just dots, you can even use iter.Seq2's in Section iteration. But brackets can expand into function arguments to do things like

	(^Type[.]^) which could call a Type function from somewhere in the context and pass it the top of the context.
	(^=Type[.]^) or include a partial named by the result of that function

This was a pretty common ask from out in the interwebs about Mustache. How can you render based on the value of some property? With Happy, it's pretty trivial. Mind bending maybe, but trivially mind bending. 

Check out the Peg template - grammar.Rules returns an iter.Seq2 so that I can control the order that it iterates over the rules. All those partials at the start of the template define the way to render expressions based on the struct name, which is then referenced by the exact Type expression in the code above. The @ even makes an appearance, leveraging the zero index to omit the first comma in a list. 

I'll add in the 6th tag at some point. It'll let you do a double render among other things, which is always something super painful to do in most template languages. It's also something that is rarely the right solution to a rendering conundrum, so I'm not in any rush. Prove me wrong.

And I really do hope that seeing all these fun little ascii emojis in your template genuinely makes you happy while you write them. I've already started naming them things like (^.^) kitty and (^/^) lefty and (^=^) rosie. OH oh oh, I totally forgot to mention - end tags are comments too.

	(^*Bob^).........(^/'s your uncle^) is totally fine. Comments/End tags can be multiline and contain anything other than a ^) or ^ )

Truly, I hope Happy works, for both your templating needs, and your soul.

### Wait, When?

So... one more thing. *testing.T gets super old after a while. Parameterized tests are great and all, but this gets confusing a little too quickly.

    func TestGraphemeNext(t *testing.T) {
        tests := []struct {
            name string
            g    *Grapheme
            want *Grapheme
        }{
            {"Initial Next", &Grapheme{"a", "bc", 1, 1, 1, 2320992, 16}, &Grapheme{"b", "c", 1, 2, 2, 2320992, 16}},
            {"New Line Next", &Grapheme{"\n", "abc", 1, 1, 1, 2320992, 14}, &Grapheme{"a", "bc", 2, 1, 2, 2320992, 16}},
        }
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                if got := tt.g.Next(); !reflect.DeepEqual(got, tt.want) {
                    t.Errorf("Grapheme.Next() = %v, want %v", got, tt.want)
                }
            })
        }
    }

In a test file full of slight variations to this, it just isn't possible to figure out what functionality is being tested, and what the test case is 
covering. It also gets complicated to make standard assertions on results, setup and teardown doesn't have a clear place to happen, yeah, just not a fan.

But it's nice to run "go test" whenever, and to have it integrated into the IDE.

I'm sure there are loads of different test patterns out there. But I am trying to learn go, not get lost in the weeds researching competing testing paradigms.
So I did the only logical thing and created my own.

When is a fairly useful take on an expectation library. Instead of parameterizing tests, you just express the variations. For instance, the above test is now.

    func TestGraphemeNext(t *testing.T) {
        next := func(g *parser.Grapheme) when.WhenOp[*parser.Grapheme] {
            return func() *parser.Grapheme {
                return g.Next()
            }
        }
        when.YouDo("Initial Next", next(parser.NewTestGrapheme("a", "bc", 1, 1, 1, 2320992, 16))).
            Expect(t, parser.NewTestGrapheme("b", "c", 1, 2, 2, 2320992, 16))
        when.YouDo("Normal Next", next(parser.NewTestGrapheme("a", "bc", 3, 17, 41, 2320992, 16))).
            Expect(t, parser.NewTestGrapheme("b", "c", 3, 18, 42, 2320992, 16))
        when.YouDo("End Next", next(parser.NewTestGrapheme("a", "", 7, 1, 53, 2320992, 16))).
            Expect(t, parser.NewTestGrapheme("", "", 7, 2, 54, -1, 0))
        when.YouDo("After End Next", next(parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))).
            Expect(t, parser.NewTestGrapheme("", "", 5, 8, 14, -1, 0))
        when.YouDo("New Line Next", next(parser.NewTestGrapheme("\n", "abc", 1, 1, 1, 2320992, 14))).
            Expect(t, parser.NewTestGrapheme("a", "bc", 2, 1, 2, 2320992, 16))
    }

Now at a glance I can see we're trying to exercise "Grapheme.Next()". The inputs are clear (the "&Grapheme" became "parser.NewTestGrapheme" because 
I didn't want to make private fields public, and I also wanted to be in a "parser_test" package). The expectations are clear. And that holds for even more
complicated tests like.

	resolve := func(key string, context string) when.WhenOpOk[any] {
		return func() (any, bool) {
			key := when.YouErr(happy.ParserFrom()("KeyName", key)).ExpectSuccess(t)
			contextJson := when.YouErr(sample.ParseJson(context)).ExpectSuccess(t)
			context := happy.ContextOf(contextJson.([]any)...)
			return key.(happy.Key).Resolve(context)
		}
	}
	when.YouDoOk("String key", resolve(`"name"`, `[]`)).ExpectMatch(t, MatchJson(`"name"`))

The api is pretty straightforward, in my opinion.

    // t is a *testing.T, doSomething() returns a thing, and something, assuming all went well, equals the thing.
    when.You(doSomething()).Expect(t, something) 

There are variations on "You", namely "YouErr" and "YouOk" that take a value and an error or boolean respectively. And there are a few variations on "Expect".

    when.You(doSomething()).ExpectSuccess(t)    // fails the test if the result is zero for its type
    when.YouOk(getIfExists()).ExpectSuccess(t)  // fails the test if the boolean return is false
    when.YouErr(doOrDie()).ExpectSuccess(t)     // fails the test if the error is not nil

There's also an "ExpectFailure(t)" that just inverts the Success criteria, an "ExpectError(t, message)" that checks that there is in fact an error returned
and that the ".Error()" on that error matches the message (note this only works for when.YouErr). And finally just in case you want to check for something more
interesting than success or equals, there's.

    when.YouOk(op()).ExpectMatch(t, SomeMatcher())

where SomeMatcher() here is an example of a when.Matcher. You'll almost certainly want to make use of the when.Assert* assertions in your matchers. See the
peg_render_test.go in happy/sample for an interesting matcher implementation.

One of the great things about You().Expect() is that it returns the value. The "key :=" line in the resolve function above grabs a parsed key. Failures don't
abort the test, so you will always get the returned value, though it may not always be useful.

Finally, I found that I spent a lot of time wrapping these You().Expect() calls in a t.Run, so that's where the "when.YouDo" calls come in. They take an executor function
of the type "when.WhenOp" as well as a name for the test case. They set up the run and execute the op within that run. So anything you might be getting from the test, for 
instance execution times, will actually be tracking the time for the operation, not just the time for the assertion.

Anyway, if it proves (potentially) useful to folks, I'll document it further. But it was a pleasant diversion before trying to tackle the next thing on my list, which, just
for the record, is no longer an Earley parser. I have a more ridiculous direction in mind.
