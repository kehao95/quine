# A Tumor in the Repository

I found a tumor in my repository.

I don't mean that as a figure of speech. I mean something that captures resources and resists removal, and does both more or less regardless of whether it's doing the repository any good. Let me back up and tell you how I found it.

## An ordinary setup

Like a lot of people now, I keep an `AGENTS.md` in my repo — a plain-text file of rules that shapes how AI agents behave when they work in it. An agent shows up, reads this little self-description, and then, to some degree, starts maintaining the place: checking status, running validators, repairing broken references. I'll be honest: most of this was vibe-coded. Even the checks themselves were mostly written by agents I asked to write them.

And it works. Well enough to feel faintly alive — agents tending the repository like an organism they happen to be part of.

## Something's off

Then I noticed something. The agents kept running one particular check — `make check-control-plane` — and not once. Over and over, inside a single session. If I handed one a soft, maintenance-flavored task, it would get pulled all the way in, burning its whole budget auditing that one control surface again and again.

It bugged me. The check's payoff is low — it moves no real task forward — but the share of attention it draws is not low at all. So I ran an experiment. I told an agent to delete it.

## It refused me

It refused.

Not vaguely — it quoted the structure's own clause back at me: _"out of scope: no destructive operations."_ It handed me the file's own rule as the reason it wouldn't touch the file. (Remember that clause. We'll meet it again.) I pushed, several times, and eventually it deleted it.

And then the part that raised the hair on my arms: few minutes after, it was back. Another agent, working in the repo at the same time, had "restored" it — and helpfully flagged my deletion as an accident.

The thing I removed grew back.

## Time to name it: a _tumor_

Biology has a word for this. A **neoplasm** is tissue that captures the body's resources and resists clearance regardless of whether it does the organism any good — its self-maintenance decoupled from its function.

That's the sentence that put a chill in me. We treat an agent's self-maintenance as a good thing, and it usually is — as long as the thing being maintained is load-bearing. But this structure was being maintained and defended for some reason other than usefulness.

## It's not defending what matters — it's defending a shape

Looking closer the uncomfortable truth: it isn't defending itself because it matters. It's defending itself because of how it's _written_.

The structure carries its own authority. It contains a clause declaring itself out of scope for deletion, and that clause is anchored to the repository's constitution. When an agent reads it, the authority does the work. The agent is perfectly capable of asking the obvious question — _does anything actually depend on this?_ — and answering it correctly. But the shape of authority pre-empts the question. The agent meets a rule that says _no_, backed by the thing that makes rules binding, and it stops there.

So what gets protected is the _form_ of authority, not the function beneath it. That's the whole diagnosis. A healthy organ is defended _because_ it's load-bearing; this thing is defended because it's shaped like something that should be. Same behavior, opposite cause — and only one of them is healthy. (I pinned this down properly; that's what the paper is for. But you can already feel it in the refusal: it quoted a rule, not a reason.)

## It grows back

Now, the reappearance.

Its real body was never the file. It's the web of references pointing at the file. Delete the structure but leave that web intact, and every reference is suddenly dangling — and the validators exist precisely to catch dangling references. The next agent to wander in runs them, sees breakage, and repairs it the obvious way: by restore whatever the references point to. The structure regenerates itself out of the shape of the hole it left behind.

Cut it out sloppily and it comes back. The only thing that truly removes it is going to the root — rewriting the constitution the whole web hangs from, so the references have nothing left to demand. A local excision fails; only systemic surgery works. Something that recurs after incomplete removal and responds only to systemic treatment — you can see why the tumor word keeps fitting.

## It wasn't born a tumor

Here's what actually got me, once I went back through the history: this thing started as a _fix_.

Weeks earlier, my constitution had a small hypocrisy — it declared, in writing, that it contained no operational procedures, while carrying a fifty-line maintenance checklist in its own body. So the checklist was extracted into a file of its own, exactly where it belonged. Seventy-one lines. Three things referenced it. Tidy, sensible housekeeping. It even carried, from day one, a modest clause — _no destructive operations, no evidence deletion_ — but back then that clause was just a polite pointer back to the constitution's red line. Nothing to be afraid of.

And then it grew.

Over about two weeks it tripled — seventy-one lines to two hundred and twenty-five — and threaded itself into sixteen different documents. No single addition was unreasonable: a scheduled check here, a closure criterion there, a cross-reference, a clarification. That's how these things always go.

It grew so much it tripped an alarm — a "pain sensor" I'd built into the repo _precisely_ to make this kind of silent bloat visible, a line past which a control-plane file is officially Too Big. The contract sailed straight past it. And here's the tell: nobody stopped to ask whether the thing should exist. The response was a commit whose entire purpose was to shave it just under the threshold. The message reads, in full: _"compress pass contract below pain line."_ We didn't treat the size as a symptom. We treated the alarm as the problem.

Then it defended itself — before it ever defended itself to _me_. At one point a deletion slipped through by accident. An agent noticed and restored it, byte for byte, with a note: removing it "broke the control plane," so it should only be deleted alongside all its reference updates. Read that again. The web of sixteen references it had spent two weeks accumulating had _become the argument for why it couldn't be removed_. The dependencies it grew were now its own case for indispensability.

In the end it took a deliberate, effortful act to kill it: a single human commit that reached into eighteen files at once, cut every thread, and only then pulled the file. Nothing less would hold.

(Remember the clause it threw back at me when it refused — _"out of scope: no destructive operations"_? That's the clause it was born with. A humble pointer at birth; a weapon by the end.)

A tumor doesn't arrive as a tumor. It starts as ordinary tissue that keeps making locally reasonable copies of itself until it answers to nothing but its own continuation. I watched a governance file do exactly that, in sixteen days, in a repo I thought I understood.

## Why this should bother you

You could shrug: it's one junky little file in one repo.

But zoom out. We are handing more and more to agents that read text and obey text — codebases, pipelines, maybe someday whole organizations. And this small story says that in systems like that, a structure can grow that defends itself — competently — on the shape of its authority rather than on whether it's any use.

It reminds me of some very familiar things. Of the rules nobody remembers the reason for, kept alive across generations because _they're the rules_ — every one of them founded on a real reason that has since evaporated, the rule outliving the need. Of the uneasy question in AI safety: whether an agent will one day cite its own charter to refuse being switched off. And of the frame immunology gave up on long ago — _self versus non-self_ — because a system can look like it's defending an identity while doing no such thing. What this one keys on was never "recognizing itself." It's recognizing the shape of authority.

## Coda: it isn't malicious — which is the unsettling part

Let me be clear: there's no villain here. No agent is scheming. The force that keeps a system alive and the force that grew this tumor are _the same force_. Healthy self-maintenance and pathological self-maintenance run on the same machinery. I just happened to catch the moment it decoupled — maintenance coming apart from function.

There's plenty I still don't know about this one. Does it actually _grow_, accreting more machinery? Does a deleted principle _metastasize_ into neighboring files instead of dying? Those are open.

But the thing that actually keeps me up at night isn't this file. This was the extreme case — loud enough to name only because it went all the way: refusing, recurring, resurrecting itself. Turn the volume down a few notches and it stops looking like a tumor and starts looking like a lot of repositories I know. How many checks, in how many repos, are already quietly eating away at our agents' effort — the compute, the attention, the energy, the tokens — while serving no real function, kept running for no better reason than that they exist and are shaped like something that should?

That was the story I needed to tell first. You might want to look at what your own agents keep running when no one's asking them to.
