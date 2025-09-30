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

Note from my future self: ah, silly past self. The whole point of making a parser is that you'll be able to do fantastic things you didn't think about doing at the time. The template engine and simple code template led to an interesting idea. A parser takes text and turns it into an AST. A template engine renders text given a set of instructions that help it walk an object graph, i.e. it takes an object and turns it into text. Imagine you had a function that turned that AST into the object you wanted to hand off to the rendering engine. So a function f that meant you could have T(F(P(input text))) where T is the template function, F is the transformation function, and P is the parsing function. So... 3.1) Pure functional transformation language, coming right up.

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

### Funki - Functional programming for the rest of us. (And the ridiculousness has arrived)

I love functional programming. Well, I love the idea of it. I was a math major, so maybe that accounts for a bit of it. But if you've ever tried to pick up Haskell or Lisp, you probably promptly put it right back down. The concepts definitely take some getting used to, and they're clever as all get out, but it genuinely takes me way too long to grok even simple functions. Reading functional code takes me back to rationalizing disassembler output. It just takes a long time to read something you didn't write.

But that's not what functions are in math. They are exist to capture cryptic ideas in clear ways, not the other way around. The beauty isn't that you can write as few symbols as possible, or obfuscate every single intention. The beauty is the idea expressed, not the mental gymnastics it takes to grok the expression.

Before we get into the functions themselves, let's talk data types. There are really just 3, although it will look like 5. Scalars are just representations for numbers.

	1     # the number 1
	1.1   # 10% more than the number 1
	-1    # the number 1, but the "other" direction
	1e1	  # 10 for those folks who like obfuscation for the sake of it
	010   # 8 for those folks who hate thumbs
	0xF	  # 15 for those folks with 16 fingers
    true  # also just the number 1, most of the time
	false # zero, definitely always
	nil   # also zero, sort of
	nan   # the number which is not a number
	inf	  # the number which is greater than all other numbers

Tuples are going to show up much more in funki than in other languages. You can think of them as fixed length arrays, but the elements don't have to be the same type.

	(1, 2, 3)  # a three element tuple
	(1)        # a one element tuple
    ()         # a no element tuple
	nil        # the same as (), sort of

Next we have Lists. They may look like arrays, strings, and objects but it's better to think of them as linked lists. Like tuples, elements don't have to be the same type.

	[1,2,3]               # under the covers this is really 3 nodes in a linked list
	[]                    # the empty list
	"abc"                 # this is the same as ["a", "b", "c"]    
    ""                    # the empty list
	{key: 1, bob: "hope"} # this is the same as [("key", 1), ("bob", "hope")]
    {}                    # the empty list
	nil                   # would you believe it's just [], "", or {}, kinda sorta?

This is why I said it's really 3, but looks like 5. Strings and objects come up so frequently that it's nice to have a convenient way of expressing them, even if it's hiding a bit of complexity in the implementation.

Ok, time for functions. Declaring a function is just naming a relation. A relation is just a mapping from a pattern to a production. That sounds more confusing than it is, let's look at an example

	f(x) -> 2 * x

here we see the name "f" given to the relation "(x) -> 2 * x". You might get the impression that a relation is just what some languages call a lambda expression, and that's close to true, but we'll see that's not exactly right soon. Here we have a pattern "(x)" and a production "2 * x". That pattern says we're expecting a tuple with a single element. so things like (1) or ("bob") would match, but () and (0, 0) would not. When the pattern matches, the corresponding elements are mapped to the given variables to be used in the production. so f(1) would map the 1 to x for use on the production side of the declaration, which would be 2*1 or just 2. Productions usually evaluate more functions, in this case, the * function is defined by the language to be multiplication of scalars.

I realize that was just a long winded way of saying functions are just how functions work. But they're a little bit more confusing than that. A pattern doesn't have to be a tuple. It does however have to be a single thing. Every function, always, always, always, takes exactly one argument, even if it doesn't "look" like it.

	f x -> 2 * x
	g(y) -> f y
	h a b -> a + b  #ERROR! this won't even parse as a function declaration
	q {age: years} -> years + 1

That last one is a bit intimidating. The pattern doesn't have to be a tuple. Here we are saying that q's argument is going to be one of those lists of tuple pairs where the first element in the tuple is a string. If we think of those pairs as the properties of an object, then this pattern is expecting a property named "age" and whatever the second element of that age property tuple is will be mapped to "years" for the production. So this q function must be something related to birthdays maybe.

What happens if q is called and there is no age tuple in the list, or the argument to q isn't a list of tuple pairs, or maybe it's not even a list at all. What happens if we try to do q(1)? Short answer, it will fail. The number 1 doesn't match the pattern, so trying to apply the relation named q to it is going to have undefined behavior. "Undefined" here means Kermit-style run around screaming. 

So... what can we do then? In math we would call these things piecewise functions, where we want to specify the behavior individually for different subsets of the inputs. But funki just calls them tests

	q x -> test x
		{age: years} -> years + 1
		_ -> 0

The keyword test tells funki to try each successive relation until one matches. That underscore is a wildcard; it'll match anything for x. So if x is an object-like list with an age property-tuple, then return the value side, or second element, plus one. If x doesn't match that age property pattern, then it'll definitely match the second pattern and return zero.

Ok, but if that object-like notation is just syntactic sugar, how would this work for lists? Remember, it's best to think of them as linked lists, so we really want to be able the say that we want the first node in the chain. Enter the pop operator, or pop-op.

	p [x:xs] -> len(x) + p xs

That colon separating x and xs says to take the first element from the list and map that to x, and the rest of the list (think of it as the "next" pointer) gets mapped to xs. Here we seem to be adding up the len's of all the elements of the argument list. It should also be noted that a tuple of one thing is the exact same thing as the one thing by itself, so for example

	p([x:xs]) -> len(x) + p(xs)

This function definition is the exact same as the previous declaration. It is recommended to use parentheses whenever it helps clarify the intention of the function, which is likely the vast majority of the time. No one is going to give you a cookie for saving those shift-9 and shift-0 keypresses.

So, we see how to get information out of a list. What if we want to return a list?

	p([x:xs]) -> len(x):p(xs)

Just like we had the colon to indicate pop on the pattern side of the relation, here we use the colon to indicate push on the production side. This says to evaluate the len of x, and then push that value in the front of whatever p of the rest of the xs returns. Hold on a second, what happens with an empty list? Don't we need to specify it?

Great observation, but no, not usually. Remember how nil was really vaguely defined? As a pattern, nil matches 0, (), and []. As a production, it's 0, (), or [] depending on what's needed when it's resolved. And every function gets an implicit test x ... nil -> nil added to the end of it. So if that's the behavior you want, and usually it is, then you don't have to specify anything. But let's say you had something more elaborate, then just use test

	elaborate(list) -> test list
		[x:xs] -> 2 * x + elaborate(xs)
		[] -> 3

Could someone do something silly like "elaborate()"? Yep, and that will return nil. In general functions with literal empty arguments are pretty rare, since they's just constants and we have better ways of expressing those. So "elaborate()" will generate a warning. But yeah, it's still a possibility. If you really want to prevent it, then you can do something like.

	supersafe(x) -> test x
		... # some relations
		nil -> panic("at the disco")

But that just feels really overkill. Panic, as you might imagine, is not a pure functional function. Turns out, I don't much care if you call things with side effects. I'm just not going to give you any way of defining side effects in your function, or managing them in any meaningful way. I'm also going to cache results, so if you're expecting side effects to influence evaluation, you're probably going to be very disappointed.

Now might be a great time to get some of the disappointment out of the way. If you're into functional languages, you might be pretty unhappy at this point. There's no type system in funki. That means all those errors that type systems insulate you from are back in full force. It's ok, folks write javascript all day long. Is it a mess? Sure, but it also powers the client side of the internet, so we apparently make do. And we might as well get currying out of the way. If you don't know what that is, you're in luck. Funki doesn't do currying. You do.

	f(x, y) -> ...         # do something with x and y
	g(x) -> (y) -> f(x, y) # calling g will return a relation that can be called 

In this example, the production side of f isn't important, it just matters that f takes a tuple of size 2. Currying means we want to create a new function that returns a function that "finishes" a call to f. Normally functional languages implicitly create that g for you (well, not exactly true, but the end result feels that way). Funki doesn't implicitly create phantom functions for you. But you can still make them if you want them, you just have to be explicit. If we can produce a relation, there needs to be a pattern to match it, right? Correctamundo

	h(z _->_, x) -> z(2 * x) / 2 

Here h is expecting a relation and something to pass to that relation. The underscores still mean you need a "something" to match, you're just not worried what that something is. If you wanted to produce something with h, you might use "h(g, 3)" with g being the explicit curry function we just declared a bit ago. Folks who are used to thinking with functions don't need all that extra information, but the rest of us don't have to even think. We see a single argument function declaration that produces a single argument relation. Can you take it to extremes?

	f x -> y -> z -> a -> (x + y) * (z + a)

Yup, but then all the rest of us just see a jerk. Could I even make it so that you could get rid of all those inner arrows and basically be able to write haskell-ese? Syntactically, maybe. Ethically... look, this is not a general purpose functional language. This is about transforming stuff you get from a parse into something useful. The principle "useful" is for handing off to a templating engine. If you want a general purpose language, there are tons out there. This is about solving a problem in a way that normal folks can read and write so that those normal folks can maintain it without having to force their brains into a sausage grinder every few months.

Ok, back on task. Let's see all the options we have for the pattern side of a relation. Remember that every function takes exactly one argument, so every function declaration must declare exactly one parameter pattern.

	f x -> ...            # If that parameter pattern is a name, then that's the name of the argument. It can be anything - scalar, tuple, list, or relation, 
	f(x) -> ...           # Tuples make functions look like functions. Use them if you can.
	f(_) -> ...           # This will match any one element tuple, but you won't be able to use that element.
	f(x, y) -> ...         # One parameter pattern can give us multiple parameters.
	f([xs]) -> ...         # This will only match a list, just remember objects and strings are lists too
	f([x:xs]) -> ...      # This will pop an element off the list as x, and give you the rest of the list as xs
	f([x:y:xs]) -> ...    # This pops two elements off the list. Note that this will only match if there are at least 2 elements in the list.
	f(0) -> ...			   # This will only match the scalar 0.
	f(false) -> ...			# This will only match the scalar 0.
	f(true) -> ...			# This will match any scalar that isn't 0, nil, or nan
	f("exact") -> ...     # This will only match the string "exact".
	f({x}) -> ...           # This will match any object-looking lists.
	f({key: value}) -> ... # This will match any object-looking lists with at least one tuple whose first element is "key" 
	                       # and make the second element available to the production as value.
	f(g _->_) -> ...       # Matches any relation.
	f((x, y), z, "crazy") -> ... # This might mean you're starting to think like a functional programmer.

So, to recap, a pattern is one of the following
*) a name
*) a literal scalar or string
*) a list
*) one or more pops off a list
*) an object
*) key-value pairs
*) a relation to an underscore with arguments that are
	*) an underscore
	*) a tuple of underscores
*) a tuple of patterns

Just a note, do you have to use underscores for relation patterns? For now, yes. I think patterns matching patterns is a little too meta for what the language is trying to accomplish right now. Also, just for clarity, you'll get a warning if you use a name in the pattern but not in the production. If you don't want that warning, use an underscore instead of the name.

Also, why is it "object-looking list" instead of just object? All the object stuff is just syntactic sugar. Let's look at it a little more in depth. That "{key: value}" pattern, how would we do that without the sugar?

	valueOf(obj, key) -> test obj
		[(k, v):o] -> test k = key
			false -> valueOf(o, key)
			true -> v

This isn't exactly a replacement for the pattern, as it's a function and functions aren't patterns. But this is how you could get a value for a key from a list of tuples that's masquerading as an object. But if that's get, how do we set?

	set(object, key, value) -> (key, value):object

Um, ok, sure, but what's the sugar version? There isn't one. I've tried coming up with one, and anything I've come up with either implies that assignment is a thing, which it isn't, or it implies that you can index into a list, which you can't. Totally open to suggestions, but this is what set looks like for now. Ternary functions that masquerade as operators are tricky on the best of days, looking at you ?:.

Ok, what about productions, what are those options?

	... -> 123 						# scalars
	... -> "happy"					# literal strings
	... -> x						# valid names, either from the pattern(s) of the function that lead to this production, or a global name
	... -> (a, b, c)					# literal tuples
	... -> [1,2,3]					# literal lists
	... -> {name: "bob", age: x}	# literal objects
	... -> test production			# tests the result of the production against successive relations.
		relation
	... -> x:xs						# push ops
	... -> x in xs					# containment
	... -> object._					# first property
	... -> object.prop				# first named property
	... -> object.[]prop 			# property list
	... -> "Ms. " ++ name 			# concatenation
	... -> 0 = 1					# equality operators, also !=, <, >, <=, >=
	... -> 0 + 1					# arithmetic operators, also -, *, /, %
	... -> x & y					# bit/boolean operators, also |, ^, <<, >>
	... -> -x						# prefix operators, also !
	... -> len(pi)					# external functions and constants (these have to be explicitly registered with the interpreter, and are not reserved)

So only 2 keywords, "test" and "in"? Oh, and the scalars "true", "false", "nil", "nan", and "inf". Well, almost. Remember way back when we were first discussing that functions always had to take exactly one argument, we said that there's a better way to define constants than using empty tuples?

	const worldsGreatestBob <- "Ross"
	bestBob(person) -> person.surname = worldsGreatestBob	

Constants are defined in a way that makes them feel a bit like a function declaration, but that would imply that the constant's name is actually an argument which all gets quite confusing. That's why the arrow points the other way. It's just not a function declaration. It's a different thing, another way of contributing a name to the global namespace.

I think there's just a few operations left to discuss. The "in" operator checks containment, like 1 in [1,2,3] or "cred" in "incredible". Concatenation gets its own operator so it's clear that 1 + 2 = 3 and "1" ++ "2" == "12". It also works for any other lists, like [2,3,5] ++ [8,13,21] or {name: "Hubert"} ++ {age: 104}. All the standard math operators show up.

Those last three are more syntactic sugar for object-like lists. Remember, they're just linked lists, so there's nothing preventing the same key from showing up multiple times.

	const guys <- {fullNames: false, name: "Tom", name: "Dick", name: "Harry"}
	... -> guys._ 			# false, it returns the second element from the first named tuple.
	... -> guys.name		# "Tom", it returns the second element from the first tuple whose first element is "name".
	... -> guys.[]name		# ["Tom", "Dick", "Harry"], hopefully obvious why.

These are all implementable within the existing language, but they're operations used so frequently in parsing that they felt like they needed to be part of the formal language.

Ok, I think that's about it. Leave your questions/comments somewhere I can find them. Probably voicemail given the folks who will even know this exists.





