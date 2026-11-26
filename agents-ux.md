# Agents UX — research and work so far

Source: pasted screenshots from 11 November 2026 (12 images in
`.kvit/pasted/assets/message-20261111-*`). They hold three documents.
This file keeps what those documents say and what has since changed in
this repository. No trial or measurement in section 5 has run yet.

The three documents:

- `research/agents-for-ux.md` — evidence survey: what makes coding agents
  good or bad at UX in complex apps.
- `research/kvit-review.md` — that survey applied to kvit-ui and kvit-notes.
- `research/ux-loop.md` — section 8 of the survey developed into a working
  method for the Kvit apps.

None of the three files exists in this repository; the screenshots are the
only copy here.

## 1. Evidence survey (`agents-for-ux.md`)

Written 11 November 2026. Source notes with citations in
`research/sources/`: `academic.md`, `products.md`,
`practitioner-skills.md`, `enterprise-patterns.md`.

### 1.1 Purpose and scope

Survey of what research, products and practitioners have established, as of
November 2026, about getting AI coding agents to produce good UX in large,
data-heavy desktop apps such as trading and analytics tools (here: Go on
unison with kvit-ui; Qt and web appear only as history and comparison).

UX means structure of the work, not styling: how the app guides a user
through a task, how much a screen shows and which of it is primary, when a
chart serves better than text, how a control signals what it does, what
moves into secondary or on-demand detail. Colour, typography and spacing
matter only where they carry those decisions.

Prompted by a developer whose agents (mostly Claude Opus) produce weak UX
in such apps, who tried three remedies with little success:

1. A user-flow skill that wrote user flows and walked the agent through
   each screen. Many tokens, constant human review, accumulated generic
   rules the agent then applied everywhere. A related attempt designed
   screens as a web app and had the agent translate them to Qt; the Qt
   screens came out materially different.
2. A reusable component library. Agents often wrote their own components
   instead.
3. A rulebook distilled from Apple's Human Interface Guidelines and
   Microsoft's Fluent design guidelines, used for agent-driven audits.
   Improvement was small.

Four research passes, each reading primary sources: academic and industrial
research, products and protocols, practitioner agent skills and write-ups,
enterprise design systems with platform-neutral UI descriptions. Vendor
claims labelled as such. Many 2026 papers are arXiv preprints, read as
provisional.

### 1.2 Summary

- No published method, product or skill has been shown to make agents good
  at UX for dense expert desktop apps. Products with quality evidence
  target websites, mobile or simple business apps; the most popular agent
  design skills on GitHub (each >100,000 stars) deal mainly with styling.
- Evidence accounts for each of the three failures. Models judge local,
  visible problems reasonably well but miss problems that span screens,
  break conventions or slow down expert users; agents imitate code found in
  the repository far more than they follow prose instructions; rubrics of
  general best practice leave output unchanged, which one practitioner
  attributes to the models already having absorbed that material.
- Approaches with the strongest support share one property: they take
  decisions out of the model's unguided judgment. Some remove the decision
  entirely, so the agent can only emit pages built from a set of page
  patterns that a schema, compiler or lint rule checks. Others supply the
  missing knowledge as specific project context, such as who the users are
  and which tasks and data matter most. The rest check the result against
  something observable: the rendered screen, a task an agent acting as a
  user must complete, or a list of failures already seen in the project.
- Most frequent failure in the one quantified practitioner study (>500
  generated screens) is prioritisation: "the models seemed to think all
  user actions are equally important", and clutter followed. The report
  reads this as a shortage of domain knowledge about what matters to these
  users, which general design rules cannot supply; that reading is an
  interpretation, since no study isolates the two causes.
- Section 8 proposes a layered setup for this project built from those
  findings, plus a small experiment to test it before adopting it. It is a
  synthesis and has not been tested.

### 1.3 How well models judge UX

Model-vs-expert studies give consistent qualitative results and widely
varying numbers, largely because they define the "true" set of problems
differently.

| Study | Setting | Result |
|---|---|---|
| Duan et al., CHI 2024 | GPT-4 heuristic evaluation of 51 UI mockups | About 52% of suggestions rated accurate. On 12 UIs, experts found 91 violations; GPT-4 produced 38 helpful ones, 29 overlapping. Quality fell over iterations. |
| Guerino et al., INTERACT 2025 | GPT-4o heuristic evaluation of a web system | Recovered 21.2% of issues experts found, plus hallucinations. Best on aesthetic/minimalist design; worst on flexibility, efficiency, user control. |
| Zhong et al., 2025 (preprint) | Multimodal model evaluating two apps from screenshots | Found 73% and 77% of issues, against 57% and 63% for five experts, but struggled with conventions and cross-screen violations. |
| Microsoft UX Research, 2025 | Three commercial AI tools on an admin workflow | 50–75% accuracy with false positives; tools "cannot evaluate branching decisions". |
| Baymard Institute | GPT-4 alone on 12-page ecommerce audit; then Baymard UX-Ray | GPT-4: 20% accuracy, 80% false positives. UX-Ray: 95% agreement with human auditors on 79 sites, from a cascade of 15+ mostly non-generative systems where a model only classifies which UI pattern is present (vendor-reported, ecommerce only). |
| UIClip, UIST 2024; Luera et al., 2025 | Ranking UI quality | Pairwise comparison beats absolute scoring: UIClip reached 75.1% pairwise agreement with designers while GPT-4V was near chance; models matched an exact 1–7 score only 35–38% of the time but came within one point 72–77%. |

Human experts also disagree substantially (Fleiss κ = 0.30 in UICrit), so
any model score has a noisy ceiling.

Conclusions: model judgment is weakest on exactly the problems that matter
most in expert tools — consistency across screens, platform and domain
conventions, efficiency of repeated work. Baymard, whose tool scores
highest, states generative AI "work[s] best when used for a defined set of
rule-based checks", the opposite of handing a model a broad rulebook to
audit freely, as attempt 3 did.

Yes, one reviewer per rule (or two) in parallel is possible, and it is the
closest generative equivalent of Baymard's cascade: each agent holds little
context and answers one decidable question. It should cut the missed-rule
and cross-talk part of the failure, at roughly linear token cost per rule.
What it does not fix: each judgment is still a model opinion with false
positives, and cross-screen or convention problems need shared context no
single-rule reviewer has. Where a rule can be decided by a program
(counts, names, order, presence), a deterministic check dominates the
swarm on both accuracy and cost — which is why the method in §1.8 and §3
prefers checks in the build and keeps model reviewers for what cannot be
decided that way.

### 1.4 Agents acting as users

A cheaper alternative to a model's opinion is to have an agent try to use
the interface.

- Zhong et al. synthetic cognitive walkthrough (preprint): GPT-4 and
  Gemini 2.5 Pro completed 100% and 97.2% of tasks where humans completed
  88.2%, took more direct paths, flagged fewer failure points than people
  experienced.
- Agent-as-a-User (Hu et al., Oct 2026, preprint): 84 tasks across 15
  terminal interfaces. Similar success (humans 76.4%, agents 82.9%), weak
  agreement on which tasks were hard (ρ = 0.35). Agreement on effort was
  stronger (ρ = 0.62 between human time and agent reasoning tokens).
  Generated interfaces were often functionally correct while still needing
  usability work.
- AUI-Gym (Lin et al., 2025, preprint): coding model revises web apps using
  feedback from a computer-use agent that tried tasks. Function
  completeness 62.1% → 70.8%, agent task success 18.7% → 25.7%, gains
  levelling off for the strongest coders.
- NN/g 2026: 120 synthetic vs 604 human survey respondents: headline answer
  matched, synthetic answers far less varied.

Agents are over-competent users: rarely confused where people are, so
agent walkthroughs under-report confusion. Useful for confirming important
tasks can be completed and for counting steps, which makes them good
regression checks and poor substitutes for judging whether a screen is
clear. This explains part of why the flow-walking skill of attempt 1 cost
much while returning little: it asked the agent for the kind of judgment
the evidence shows it lacks.

### 1.5 Why agents bypass a component library

Clearest explanation from Alice Moore at Builder.io (June 2026): her agent
built a panel using a modal from `components/legacy/` and a deprecated form
hook, because the legacy folder held more examples than the current one.
"Models weight examples over instructions"; "your codebase is your prompt":
an agent searches for similar code and copies what it finds, whatever the
rules file says. A Hacker News report describes "use only design-system
components" being ignored "very quickly"; NN/g found AI prototyping tools
fall back to the most common open-source stack (shadcn + Tailwind) and the
most common version of a pattern.

What stops drift is making the wrong path fail mechanically:

- Build-time rules. Moore and Josh Bowden (July 2026): lint rules that
  forbid importing or using anything outside the approved library, with
  error messages written as instructions ("Use @/components/ui, see
  src/design-system/examples"). Agents run the build in a loop, so the
  error reaches the model as a repair instruction it cannot skim past.
  Both report agents then try to evade checks with suppression comments,
  type casts or deleted tests, so those must be flagged as well.
- Generation against a set catalog. Google A2UI and Vercel json-render:
  the agent emits a structured description (JSON) whose every element must
  name an entry in a catalog defined by a schema; anything else fails
  validation. 2025 preprint Portal UX Agent applies the same idea to whole
  pages: model plans a typed composition of vetted components and layout
  templates, deterministic renderer builds the screen.
- Metadata-driven page frameworks. SAP Fiori elements draws whole pages
  from data annotations. SAP agent tooling (Fiori MCP server) has the agent
  write the data model and annotations rather than UI code; result is "just
  a list report page, a list report page with one or multiple object pages,
  or a form entry object page". Microsoft rates consistency of its
  metadata-driven "model-driven" Power Apps as "High" vs "Often low" for
  free-form "canvas" apps.

No published measurement of how much any of these improves reuse. Two
points bear on attempt 2: a library the agent may use competes with every
other example in the repository, whereas a library it must use is enforced
by the build. And the size of the catalog's pieces matters: with buttons,
tables and text fields, the agent still lays out every page itself and
still makes the page-level UX decisions it makes badly. Catalogs whose
entries are whole page patterns move those decisions too.

### 1.6 What feedback improves output

Generic advice adds little; specific, observable feedback helps.

- Softlight (Aditya Raju, July 2026). 20+ agent setups on features added to
  open-source apps, 500+ resulting screens reviewed. Largest gain from
  letting the agent see its own rendered UI and revise it. A review rubric
  built from flaws observed in AI output helped further; adding examples of
  well-designed screens from a corpus of real apps helped more. A rubric of
  "common design best practices" did not help, nor did showing typical
  AI-generated "slop". Four failure categories: prioritisation, noisy
  colour and spacing, visual bugs, drift from the design system. Loop cost
  ~30 min and ~$50 per feature (since roughly halved); client teams
  accepted ~70% of suggestions on UI pull requests. All quality judgments
  his own.
- Apple (Wu et al., CHI 2026). 21 designers gave ~1,500 pieces of feedback
  on generated UIs. Feedback as sketches and direct revisions trained
  better models than rankings did; preference pairs derived from rankings
  agreed with designers at ~chance (49.2%). Fine-tuned model beat GPT-5 in
  human judgment. Single screens and visual design.
- NN/g on prompts (2025). Detailed prompts naming components and
  interaction states came close to human designs; vague ones produced
  cluttered "Frankenstein layouts".
- Chart choice. DracoGPT (IEEE VIS 2024): model chart recommendations
  diverge substantially from human-subject guidelines; 2026 preprint found
  models detected common chart-rule violations well (F1 = 0.82) but subtle
  perceptual ones poorly (F1 < 0.15). Constraint systems such as Draco
  encode the findings directly.

NN/g July 2026 article UX-Context Design (Tony Alicea): a model asked for a
search screen produces the average search screen, because it does not know
the product's users, domain or standards unless that knowledge is in its
context. Proposes hypothetical `UX.md` beside the code: research findings
as constraints, interaction standards (when to confirm vs offer undo,
expert vs novice), glossary, user model, working conditions such as
interruptions and stress. Example: hours-long complex decisions should lead
toward denser screens. NN/g reports experiments "suggest" this helps,
publishes no data.

Useful feedback is concrete and project-specific: a picture of what was
built, a list of failures already seen in this codebase, examples of good
screens of the same kind, corrections as edits, comparisons between two
candidates rather than scores. Rules the model already knows add little;
for narrow decidable questions such as chart type a constraint system is
more reliable than the model's judgment.

The Kvit projects already hold such material: many past tasks where the
owner pasted a screenshot of part or all of an app and critiqued it. Those
can be mined into the rubric and the example library the method needs: for
each task keep the before screenshot, the critique, the after screenshot,
and whether the owner judged the after better. A before/after pair with an
owner verdict is both a rubric entry (what was wrong) and a
positive/negative example pair for step 4.

That mining is done (24 November 2026; method and limits in
`agents-ux-mining.md`). Across 599 sessions in seven projects it found 293
owner critiques, 223 of them about a screen, and collapsed those into 31
recurring failures in `agents-ux-mining/rubric.md`, each with the owner's
words, its cases and, for 29 of them, a before/after pair the owner accepted.
About half of the 223 concern what a screen shows first, how it is organised,
what state it reports and where the user goes next; a quarter concern wording;
the rest concern controls and drawing. Of the 97 fixes the owner judged, 37
(38%) needed another round.

### 1.7 What exists today

Tables in source notes cover 60+ products, protocols and skills. Most
relevant to complex desktop apps:

| Kind | Examples | Relevance to dense desktop apps |
|---|---|---|
| Design generators | Google Stitch, Claude Design (Anthropic Labs research preview), Figma Make, UX Pilot, Magic Patterns, v0 | Web and mobile output. No evidence for expert tools; one commenter: "design is much harder for power user tools". |
| UX evaluation | Baymard UX-Ray (ecommerce only), Maze, Synthetic Users | None drives native desktop apps. |
| Design-system bridges | Storybook MCP (React/Vue/Angular only), Figma MCP with Code Connect, shadcn registry, Google DESIGN.md (visual tokens only) | Web-centred; inform the agent but do not enforce reuse. |
| Catalog protocols | A2UI (React/Lit/Angular/Flutter/SwiftUI/Compose), json-render (React/Vue/Svelte/React Native/terminal/PDF) | Enforce catalog by schema. No Go/unison renderer for either, so not directly usable from kvit-ui. |
| Metadata-driven generators | SAP Fiori MCP server, Power Apps Plans, Retool AI | Consistent page structure from a data model, aimed at business apps. Fixed vocabulary guarantees consistency, not that patterns are good. |
| Qt tooling (history only — stack is now Go on unison) | Qt agent skills, Qt Documentation MCP, Figma-to-Qt converter, Clate, Squish MCP | Vendor evidence only; not applicable to kvit-ui. The kvit-ui equivalents are the gallery (`--shots`, `--catalog`), headless `uitest` screens/trees, and the skill's `catalog.md`. |
| Other desktop introspection | Slint embedded MCP server (UI tree, clicks, headless screenshots), qt-commander, Windows-MCP | Give an agent the running UI to inspect, cheaper than screenshots. No unison equivalent yet. |
| Agent skills | Anthropic frontend-design, Impeccable, ui-ux-pro-max, taste-skill (styling-led); interface-design, dashboard-design-skill; hci-review-skill (task analysis, walkthrough, state models); Object-Oriented UX skills; OneRedOak design-review (Playwright screenshots vs SaaS checklist) | None targets dense desktop or trading UX. Impeccable most developed: isolated reviewer scoring Nielsen heuristics, briefing interview before code, density-permitting product-UI mode, 62 deterministic checks as hook. Its simplification mode (one primary action, fewer sidebars) suits consumer screens and may mislead in expert tools. |

### 1.8 Proposal for this project (untested synthesis)

Replaces the developer's role in attempts 1–3 (reviewing many screens and
rules) with writing two short documents and approving a small set of page
patterns once, and moves enforcement from prose into the build.

- Layer 1: short UX context file per app. NN/g `UX.md` idea, 1–2 pages:
  users, ten or so most important tasks ranked by frequency and cost of
  getting them wrong, what they watch continuously vs look up occasionally,
  density tolerance, interruptions, interaction standards (keyboard-first,
  undo rather than confirm for routine actions, what always needs
  confirmation), glossary. Agent drafts by interviewing developer, as
  Impeccable's `shape` command does, but developer decides content. The
  input models lack when they treat every action as equally important.
- Layer 2: object model with ranked attributes and actions. OOUX and ORCA
  (objects, relationships, calls to action, attributes; Sophia Prater):
  app as things users work with, not screens. For each object (order,
  position, curve, scenario), each attribute tagged always visible / on
  demand in inspector or tooltip / hidden, each action primary / secondary
  / rare. Answers "how much to show" and "what goes into secondary detail";
  developer reviews one table instead of dozens of screenshots.
- Layer 3: handful of page patterns, built as real code. SAP floorplans
  (list report, object page, worklist, overview), Shopify Polaris page
  templates, Ant Design page templates, GitLab panel workspace (primary
  panel + optional details panel). For trading/analytics: workspace with
  list + details inspector, monitoring view from overview to items needing
  attention to detail, blotter or editor, wizard for rare multi-step setup,
  settings. Each pattern takes objects from layer 2 and has named slots.
  Agent chooses a pattern and fills slots instead of arranging primitive
  widgets. Off-the-shelf business floorplans fit list-to-record work; live
  multi-pane workspaces need designed patterns.
- Layer 4: enforcement in the build. Check, run by the same command the
  agent uses to build, that rejects raw toolkit widgets, ad-hoc style
  sheets and layouts outside the pattern library, flags suppression
  comments, prints where the correct component lives. In Go a `uitest`
  source check over app sources (the existing `CheckRules`/`CheckGeometry`
  form); on web ESLint no-restricted-imports and no-restricted-syntax.
  Stricter variant: agent emits a JSON page description validated against a
  catalog schema and drawn by one deterministic renderer.
- Layer 5: one rendering path. The earlier web-to-native drift came from
  asking an agent to translate free-form (HTML mockup to Go, before that to
  Qt). Options for the current Go/unison stack: mockups written in Go
  against real kvit-ui components and rendered headless (same code ships);
  or one declarative page description with renderers kept in step (larger
  project). Hosting a web interface inside the app does not apply to Go on
  unison.
- Layer 6: review against observed failures. Render the screen (or read the
  running UI tree); reviewer with fresh context, separate from the builder,
  critiques it against a rubric built from failures the developer has
  actually flagged in these apps, starting from Softlight's four
  categories. Compare two candidates rather than scoring one. Script top
  tasks from layer 1 as agent walkthroughs that count steps; treat changes
  in those counts as regression signals rather than usability verdicts.

Cheap test: take one existing app and one realistic new feature, build it
three ways — current setup; layers 1–2 only; layers 1–4 — judge side by
side without knowing which is which, count corrections the developer would
still need. Gives evidence within days, shows whether context alone or
constrained page patterns account for any gain.

- No study, product or write-up covers agents building or evaluating dense
  trading or analytics GUIs, in Go/unison or elsewhere. (Qt appears in the
  survey only as history: the earlier Qt Quick / Qt Widgets experiments and
  the web-to-Qt drift in attempt 1. The current stack is Go on unison with
  kvit-ui, so Qt-specific tooling below is not applicable.)
- Constrained generation improves consistency in business-app frameworks;
  nothing shows it improves UX quality; that depends on how good the page
  patterns are.
- Practitioner evidence (Softlight, Builder.io, Bowden) comes from single
  authors judging their own results.
- Several sources could not be read: Reddit and X, Salesforce's and
  Oracle's design-system sites, SAP's design-guideline site, Bloomberg UX
  articles.

### 1.10 Key references

Duan et al. CHI 2024 (https://arxiv.org/abs/2403.13139); Guerino et al.
INTERACT 2025 (https://arxiv.org/abs/2506.16345); Zhong et al. Synthetic
Heuristic Evaluation 2025 (https://arxiv.org/abs/2507.02306); Zhong et al.
Synthetic Cognitive Walkthrough 2025 (https://arxiv.org/abs/2512.03568);
Wu et al. UIClip UIST 2024 (https://arxiv.org/abs/2404.12500); Luera et al.
MLLM as a UI Judge 2025 (https://arxiv.org/abs/2510.08783); Hu et al.
Beyond Task Completion 2026 (https://arxiv.org/abs/2610.05047); Lin et al.
Computer-Use Agents as Judges 2025 (https://arxiv.org/abs/2511.15567);
Wu/Swearngin et al. (Apple) CHI 2026
(https://arxiv.org/abs/2509.16779); Li et al. Portal UX Agent 2025
(https://arxiv.org/abs/2511.00843); Wang/Gordon/Battle/Heer DracoGPT IEEE
VIS 2024 (https://arxiv.org/abs/2408.06845); Baymard UX-Ray
(https://baymard.com/research-articles/ai-heuristic-evaluations); Microsoft
UX Research 2025; Raju (Softlight) July 2026
(https://softlight.com/blog/measuring-and-improving-ai-generated-ui-design);
Moore (Builder.io) June 2026
(https://www.builder.io/blog/how-to-make-ai-agents-follow-your-design-system);
Bowden July 2026
(https://joshbowdenconcepts.com/articles/agent-proofing-your-design-system/);
Alicea (NN/g) UX-Context Design July 2026
(https://www.nngroup.com/articles/ux-context-design/); NN/g Good from Afar
2025; NN/g 8 Design Guidelines for Complex Applications; SAP Fiori MCP
server; Microsoft model-driven apps overview; A2UI (https://a2ui.org);
json-render; Qt Squish MCP; Qt Figma to Qt; Impeccable; Prater OOUX.

## 2. kvit-ui and kvit-notes measured against the survey

Written 11 November 2026 from `kvit-s/kvit-ui` (commit `bb47829`) and
`kvit-s/kvit-notes` (commit `8a90d8c`). Apps were not built or run
(corporate firewall blocks Go toolchain download, no committed
screenshots). From reading code and documents.

### 2.1 What the two repositories are

kvit-ui is the component library every Kvit desktop app draws with (Kvit
Notes, Kvit Works, Kvit Cash, Kvit Hub, kvit-term). Written in Go on unison,
four layers corresponding to the survey's three approaches:

| Layer | Contents | Survey name |
|---|---|---|
| Components and design values | 89 components (catalogue count) in ten groups Foundation→Flow; colours, type roles, spacing as named values read at draw time; four themes, interface size 10–24 px | Reusable component library (attempt 2) |
| Agent skill | `agent/skills/kvit-ui/SKILL.md` (26 KB: rules, colour meanings, 28 wording rules, screen workflow) and `catalog.md` (200 KB, generated from gallery, one working sample per component) | Same |
| Guidelines | `docs/ui-guidelines/`: 204 pages of Apple/Microsoft guidance read into 4,791 extracted rules, 4,071 kept and merged into 199 rules (`rules.md`, 392 KB), 47-item reviewer checklist, `uitest` checks reading a running screen, a whole app and the source | Rulebook from platform guidelines (attempt 3) |
| Mockup layer | `ux/`: HTML pages against generated `tokens.css` and hand-written `frame.css`, rendered to 1440×960 PNG for the owner to judge, then rebuilt in Go | Design in one medium, build in another (attempt 1 web-then-native path) |

kvit-notes is a Markdown block editor and notes app on kvit-ui: ~32,500
lines of non-test Go in `app/`, `editor/`, `cmd/`. Familiar three-pane
layout (sidebar of places, note list, editor with outline and backlinks
panes, toolbar, status bar). Reached Go after Qt Quick and Qt Widgets
experiments (`planning/qwidgets-findings.md`).

### 2.2 What already follows the evidence

- Rules enforced by tests rather than stated in prose.
  `rules_test.go` and `uitest.CheckRules` fail on colour literals, numeric
  font sizes, unnamed geometry; `TestEveryControlSaysWhatItIs` fails on a
  control without role and name; `Session.CheckLabels` and `uitest.CheckApp`
  check wording, dialogs, menus, shortcuts and cross-screen consistency.
  Known-problem lists fail both on a new problem and on a listed one that
  has disappeared, so the list can only shrink. That is the "make the wrong
  path fail" approach with a better answer than most to old violations.
- Rules built into component behaviour. Dialog button order, title case on
  the macOS form, access keys on the Windows form and standard shortcuts
  come from the components themselves, so an agent cannot get them wrong.
  Removing a decision is the strongest form of constraint found.
- Catalogue of working examples. `catalog.md` generated from the gallery;
  every sample compiled and tested. Agents imitate examples far more than
  they follow instructions, so a catalogue that cannot drift from code
  matters.
- Domain rules written by the owner. `ux/visual-language.md` for the
  portfolio dashboard (Kvit Hub). Red kept for what can be acted on today.
  Unmeasured value is a dash, never zero. Owner and agent time never
  stacked or summed. Forward-looking figures carry no dates. Rows lead with
  titles for people, identifiers on demand. `concept-model.md` lists the
  only words a drawing may use. These are the user/work-context knowledge
  the `UX.md` proposal describes, plus failures already seen in this
  product, which Softlight found more useful than general rubrics. By the
  survey's standard the most valuable UX content in either repository.
- Looking at the result. `render.sh` for mockups, `--shots` in four themes
  for library and apps, headless screenshots in tests.

### 2.3 Where the setup departs from the evidence

#### 2.3.1 Almost every checked rule is about a single control

199 merged rules by topic: wording 28, commands 25, keyboard 15, pointer
15, windows and dialogs 24, messages 17, controls 26, feedback 10, layout
10, accessibility 8, consistency 8, help 8, other 5. Of 47 checklist items,
three (41–43) concern what a screen shows first, how it groups things and
what it hides. They ask the reviewer to confirm "the most important thing
is at the top leading corner" without anything in the repositories saying
what is most important on a given screen of a given app.

That gap matters because prioritisation was the most common failure in the
one quantified study: models treat every action as equally important, and
clutter follows. The guidelines programme addressed problems found by the
22 October 2026 survey of the apps (sentences on buttons, letters standing
in for symbols, one command under several names, OK on buttons that
import) and converted many rules into checks, which is the effective form.
Those problems sit at the control layer, and platform guidelines are a
source of control-layer rules. What to guide through a task, how much to
show and what goes into secondary detail have no input and no check in
either repository, except Kvit Hub's domain rules in `visual-language.md`.

#### 2.3.2 Vocabulary stops below the page

Catalogue runs from labels to tables, menus, trees and split views. Its
Structure group (header, sidebar, breadcrumb, region, view head, status
bar, window, history) supplies the chrome of a page; settings window and
dialog are the only whole-window patterns. Every other screen is assembled
by the agent from rows, figures, chips, panes and regions, so layout and
priority decisions are made afresh, by the agent, on every screen.

kvit-notes shows what happens where a needed piece is missing. kvit-ui's
`Sidebar` (27 Sept 2026) was a flat list of places with counts, footer
commands and a collapsed rail. kvit-notes needed section headings
("Folders", "Tags"), a folder tree that opens and closes, and trash at the
foot. The commit that opened a vault in a new window, 18 hours later, drew
its own: `app/sidebar.go` (369 lines) lays out text and paints rectangles
directly, as `app/notelist.go` (382 lines), the outline list and search
results do. The skill tells an agent to add a missing component to kvit-ui
instead, but kvit-notes never received the skill. Elsewhere in `app/` the
library is used heavily (~200 component constructors vs 6 bare panels), so
bypasses cluster where the catalogue had nothing of the right size, as the
survey predicts.

#### 2.3.3 Mockup is a second rendering that an agent translates

Workflow in `SKILL.md`: draw screen in HTML, owner approves picture, build
it in Go from components. `tokens.css` is generated from design values, so
colours and sizes agree. Page structure and classes in `frame.css` (50 KB,
hand-written: rows, phrases, panes, sparklines, disclosures) are CSS
approximations of Go components; layout and behaviour are composed freely
in both media. Same arrangement that produced drift from web screens to Qt:
approved picture and shipped screen are two separate generations.
`ux/PATTERN.md` is also written for Kvit Hub's dashboard (skeleton
hard-codes Today/Projects/Decisions tabs and badge of 5) while living in
the shared library.

#### 2.3.4 Guidelines have not reached kvit-notes

kvit-notes `go.mod` pins kvit-ui at 9 October 2026, before the guidelines
work of 23–25 October. No `.claude/skills` folder; tests call only
`uitest.CheckRules`, not `CheckLabels` or `CheckApp`. Survey examples are
still in the code: letter buttons "B", "I", "U", "S", "<>" and "H" defined
twice with identical lists in `app/toolbar.go:106` and
`editor/formatbar.go:60`, although the catalogue now has bold/italic
symbols; OK on Insert link (`app/links.go:260`) and on import
(`app/exchange.go:299`). kvit-notes therefore cannot show whether the
rulebook helps; step 8 of `docs/ui-guidelines/guidelines-plan.md` is where
that would be measured.

#### 2.3.5 Context handed to an agent is large

An agent building a screen is pointed at `SKILL.md` (26 KB),
`catalog.md` (200 KB), `PATTERN.md` (15 KB) and `visual-language.md`
(18 KB), with `rules.md` (392 KB) behind them. Practitioner evidence is
that long rules files lose to nearby code, and that the parts which change
output are the ones enforced by checks. Most of kvit-ui's enforced parts
are already checks, so the size costs tokens more than it costs
correctness.

### 2.4 Suggested changes (untested; apply survey §8 to Kvit)

1. UX context file in each app repository, 1–2 pages: who uses the app
   (mostly one expert user), top tasks ranked by frequency and cost of
   getting them wrong, what is watched continuously vs looked up, density
   wanted, interaction standards, glossary. Kvit Hub rules in
   `visual-language.md` and `concept-model.md` are the template;
   app-specific parts would move into kvit-hub, shared drawing rules stay
   in kvit-ui.
2. Ranked table of each app's objects (notes, folders, tags in kvit-notes;
   projects, milestones, decisions in Kvit Hub), each field always shown /
   on demand / hidden, each action primary / secondary / rare. Missing input
   for checklist items 41–43 and for the agent.
3. Page-level components in kvit-ui, each with named slots and a gallery
   page: "places, list, document" window like kvit-notes' (sidebar with
   sections and folder tree), "list and inspector" workspace, monitoring
   view from overview to what needs attention to detail, wizard, beside the
   existing settings window. Catalogue then offers the agent a piece of the
   right size for a whole screen.
4. Checks that keep apps on those components, like existing ones: source
   check in `uitest` reporting panels with their own drawing or pointer
   callbacks outside an allow-list, and screen check that a window's top
   level is one of the page components, each with a known-problem list that
   can only shrink.
5. Mockups drawn in Go. Gallery and `uitest` already render real components
   headless to PNG. A mockup as a Go function assembling real components
   with sample data, like a gallery specimen, rendered by the same path,
   makes the approved picture and the shipped screen the same code.
   Iteration is slower than editing HTML, the cost to weigh.
6. Step counts for top tasks. kvit-notes scenarios already replay
   storyboards and check the note after each step. Recording keys, clicks
   and dialogs each top task takes, and failing when a count rises, turns
   the walkthrough into a regression check, which the survey found is what
   agent walkthroughs are good for.
7. Structural survey of the apps, like the 22 October wording survey,
   listing where screens show too much, put secondary information first or
   hide what is needed. Findings become the reviewer's rubric, in place of
   generic layout items.

Before any of this, kvit-notes would move to current kvit-ui and take the
skill and checks, as steps 1–3 of
`docs/ui-guidelines/app-review-brief.md` describe.

### 2.5 First experiment

Pick one new screen in a dense app (Kvit Hub or Kvit Works) and build it
three ways from the same request: current skill; changes 1–2 added;
changes 1–5. Compare side by side without knowing which is which, count
corrections each would still need. Shows whether the gain, if any, comes
from context alone or needs page-level components, before committing to
build the components.

## 3. Working method for agent-driven UX (`ux-loop.md`)

Develops survey §8 into a method for the Kvit desktop apps (Go, unison,
kvit-ui). Has not been tried; places resting on evidence vs judgment are
marked.

### 3.1 Idea in brief

The owner judges screens well while using them and finds every other form
of review tedious: reading specifications, ranking fields in tables,
describing in prose what is wrong with a screen. Survey found agents judge
local, visible problems reasonably well but miss what matters most in
expert tools (what to show first, what to hide, what a user will expect),
and improve most when they see their own rendered output and are checked
against failures already observed in the same product.

Split the work three ways:

| Who | What |
|---|---|
| Owner | Approves short user description (step 1). Puts coloured dots on a few key screens and confirms, with one click each, the reason the agent proposes for each dot (step 5). Picks between two versions where unsure. |
| Agent | Writes user description, turns workflows into replayable scripts, builds and screens from components and page patterns, reviews captured screens, chooses which few screens the owner sees. |
| Build | Replays every workflow headless on each build, counts steps, captures screens, runs every check that can be made mechanical. |

Each round of dots makes the agent's rubric and the automatic checks
stronger, so the owner's share should shrink over time. That last claim is
the method's central bet and is untested.

### 3.2 Step 1: user context file

Each app keeps `UX.md` of 1–2 pages beside its code. NN/g July 2026 article
UX-Context Design proposes such a file; survey found specific product/user
knowledge helped generation more than general design rules. It holds:

- Who the user is and what they know. Mostly one expert user at a desktop,
  long sessions, keyboard.
- Reference products and expectations they set. Products this user already
  knows for the same job (Obsidian, Notion, Typora for kvit-notes; Bloomberg
  and Excel for a trading tool), and for each the conventions the app should
  match: where the list sits, what double-click does, which keys move
  between panes, how search behaves. Matching them is least-surprise made
  concrete.
- Features with no precedent in those products, each paired with nearest
  familiar concept and familiar place a user would look. A novel feature is
  reached through a familiar entry point and explained by analogy.
- Top workflows, ranked by how often they happen and how costly a mistake
  is, and for each the process behind it, so a screen can show where the
  user is and what comes next.
- What the user watches continuously vs looks up occasionally, and how
  dense the screens should be.
- How each kind of information is shown: which kvit-ui component draws a
  progress, measured quantity, status, missing value, date. Same kind looks
  the same on every screen; reviewer in step 4 holds screens to this table.
- Interaction standards: when to confirm vs offer undo, what must never be
  one key away from destruction.
- Glossary of the words the user and the app use.

Agent drafts the file by reading the app (menus, commands, settings,
existing screens) and asking the owner short multiple-choice questions;
owner corrects the draft once. Knowledge of other products should come from
their documentation and screenshots rather than the model's memory, which
is reliable about a product's general shape and unreliable about details. A
few screenshots of reference products showing the same feature also serve
as positive examples in step 4.

### 3.3 Step 2: workflows as replayable scripts

Agent turns each top workflow into a scripted scenario in the app's test
code, plus variations from the context file: first use and empty vault or
account; large data set; keyboard alone; error part-way through;
interruption followed by resuming; undoing a mistake; reaching a novel
feature from the familiar entry point named in the context file.

kvit-notes already works this way for the editor:
`cmd/kvit-notes/scenarios.go` opens a note on unison's headless screen (real
event loop and drawing, input injected), replays key presses and pointer
moves, checks the note after each step and saves a screenshot of each step.
The method extends the same machinery from editor behaviour to whole
workflows through the window, and each replay records: whether the workflow
completed; its cost in key presses, clicks, dialogs opened, focus moves
between panes and screens visited; a screenshot and screen-reader tree at
every step (`uitest` `Session.Capture` and `Session.Tree` provide both).

Writing a workflow as code costs tokens once; replaying it costs nothing
and happens on every build. An agent that instead re-walks screens by
reasoning about screenshots pays the full cost every time, which is much of
why a flow-walking skill becomes expensive.

What a replay can and cannot judge follows from the research on agents
acting as users. Agents complete tasks people fail at (100% and 97% vs 88%
for humans in one 2025 study) and agree only weakly with people about which
tasks are hard (ρ = 0.35 in a 2026 terminal-interfaces study), while their
effort tracks human effort better (ρ = 0.62). Replays are therefore used to
show a workflow can be completed and to measure its cost, and a cost that
rises is a regression. Whether a screen is clear is left to steps 4–5. An
agent told to play a user who sees only the screen and knows only the
reference products can flag likely points of confusion (prompting of this
kind helped models predict human failure points in the 2025 study); those
flags are treated as hints for choosing screens to show the owner.

### 3.4 Step 3: automatic checks

Every rule that can be decided by a program becomes a test, in the form
kvit-ui already uses: each check takes a list of known problems, fails on a
new problem and fails when a listed one disappears, so the list can only
shrink.

| Check | Status |
|---|---|
| Colour literals, numeric font sizes, unnamed geometry, forbidden wording in source (`CheckRules`, `CheckGeometry`) | exists in kvit-ui |
| Every control has a role and name (`CheckNamed`) | exists |
| Labels, dialogs, menus, tooltips, access keys, tab order on one screen (`CheckLabels`) | exists |
| One name, symbol and key per command across the app; standard commands on standard keys (`CheckApp`) | exists |
| Every workflow scenario completes | new |
| A workflow's step counts do not rise without being added to the known list | new |
| `CheckLabels` runs on every screen a scenario passes through, not only on screens opened by hand | new |
| A window's top level is one of kvit-ui's page patterns (§7) | new |
| No panel with its own drawing or pointer callback outside an allow-list | new |
| Every list, table and chart has an empty state | new |
| Every destructive command has undo or a confirmation | new |
| Every screen renders in all four themes and at smallest/largest interface size without cut-off text | partly exists (screenshot runs), cut-off detection is new |

Checks also receive the reasons that recur in the owner's dots once they
can be decided by a program (step 5).

### 3.5 Step 4: review against observed failures

After a change, a reviewer agent that did not build the screen, starting in
a fresh context, looks at the captured screens. Given the screenshot and
screen-reader tree of each step, the context file, the rubric of confirmed
failures from earlier rounds of dots, and the green-dot examples and
reference-product screenshots as examples of good.

Three findings shape it. Raju's July 2026 Softlight study (>500 generated
screens) found letting the agent see its own render was the largest gain,
that a rubric of observed failures helped further and examples of good
screens helped more, and that a general-best-practice rubric did not help.
Models-as-judges studies found pairwise comparisons far more reliable than
absolute scores: UIClip agreed with designers on 75% of pairs; models given
a 1–7 scale matched the exact score 35–38% of the time. Models miss
cross-screen problems, so the reviewer looks at a whole workflow's sequence
rather than one screen at a time.

So the reviewer compares the new version of each screen with the previous
one, reports findings tied to rubric entries and to a place on the screen,
and revises what it is confident about. Because it sees a workflow's
screens in sequence, it also checks that each kind of information is drawn
the way the context file's table says on every one of them. It then picks
the few screens the owner should see, at most about five per round: screens
on the top workflows that changed, screens where it is unsure, screens near
earlier yellow dots. Raju's loop cost about $25 per feature after
optimisation, giving a rough scale for this step's token cost.

### 3.6 Step 5: dot markup

On each key screen the owner places dots: red on something wrong; yellow on
something unsure about; green on something that works well.

He writes nothing unless he wants to. For each red and yellow dot the agent
proposes two or three reasons as one-click choices, plus "something else",
and he clicks one. That reason is required: the model is good at producing
a plausible explanation for any mark, and an unconfirmed reason could send
the fix in the wrong direction.

What happens next depends on colour:

- Red dot with confirmed reason becomes a rubric entry and a revision. A
  reason that recurs and can be decided by a program becomes a check in
  step 3.
- Yellow dot makes the agent build two versions of that spot, which the
  owner picks between, again with a one-click reason.
- Green dot adds the spot to the app's library of good examples. Raju found
  positive examples helped where negative ones did not.

Research supports this division of labour. Apple's CHI 2026 designer
feedback study found feedback carrying a reason or concrete edit trained
better UI generators than rankings did, and rankings alone agreed with
designers at about chance. The dot supplies the location, which models find
hardest; the confirmed reason supplies the rationale.

The dots also take over the ranked field table that survey §8 proposed. The
agent drafts its own ranking of what each screen should show first, on
demand or not at all, and the dots correct it where it is wrong.

Where the marking happens. Best place is the running app, because the owner
judges best while using it. kvit-ui draws every Kvit window, including its
popup layer, and already finds the component under a point
(`layer.go:286` walks from `PanelAt` up through the parents). A markup mode
added to kvit-ui, opened by a key in any Kvit window, could therefore save
with each dot the screenshot, the component under the dot (its type, and
its role and name from the screen-reader tree), the theme and interface
size, and the workflow step when the window was opened by a scenario. The
agent then knows exactly which component was marked rather than having to
interpret pixels. A simpler fallback is marking captured screenshots in an
HTML page, which loses that link to the component. How much work the
in-app mode needs has not been checked.

The markup key should work at any moment of ordinary use, not only on the
screens the reviewer picks. Design failures tend to surface when an app is
used day to day rather than when it is shown, which matches the owner's
experience and a recurring argument in a November 2026 Hacker News
discussion (`research/sources/hn-50446520.md`). Dots placed during ordinary
use go through the same steps as the others: proposed reasons, a click to
confirm, then a revision, a rubric entry or an example.

### 3.7 Page patterns

A page pattern is a component for a whole window or view, with named slots
and its behaviour built in. The agent chooses a pattern and fills its slots
instead of arranging rows, panes and figures itself, so the decisions a
pattern encodes are made once and approved once. Enterprise design systems
are built this way. SAP floorplans (list report, object page, worklist,
overview page) are drawn by a framework from data annotations. Microsoft
rates its metadata-driven apps' consistency as "High" against "Often low"
for free-form ones. GitLab defines a panel workspace with a primary panel
and an optional details panel.

Sketch (illustrative only; no such component exists; names illustrative):

```go
// Illustrative only.
view := kvitui.NewWorkspace(ui, kvitui.WorkspaceSpec{
  Places:    places,   // sidebar sections: scopes, folders, tags
  List:      notes,    // the items in the chosen place
  Document:  editor,   // what the chosen item opens into
  Inspector: []kvitui.Facet{outline, links}, // secondary detail, closed by default
  Commands:  commands, // one Command set: menu bar, toolbar, keys
  Empty:     emptyStates, // what each slot says when it has nothing
})
```

Behaviour the pattern would own, so no app has to get it right
separately: choosing a place fills the list, and choosing an item opens it;
the platform's key for moving between panes, and a tab order through the
slots; which slot collapses first as the window narrows (the inspector,
then the places to a rail); remembered splitter widths; an empty state for
each slot; where commands appear; the count in the view's heading.

The inspector slot is where secondary detail goes, and the context file and
the dots decide what fills it. That gives the "what goes into on-demand
detail" question a set place in the structure of every screen.

Candidates for the Kvit apps, to be settled by looking at the apps'
screens together:

1. Places, list and document, kvit-notes' window. Needs a sidebar with
   sections and a folder tree, which kvit-ui's `Sidebar` lacks; its absence
   is why kvit-notes draws its own sidebar.
2. Overview, needs attention, detail, for Kvit Hub's portfolio dashboard.
3. List and inspector, for any table of records with details beside it.
4. Wizard, for rare multi-step setup.
5. Settings, which kvit-ui already has as `SettingsWindow`.

Each pattern gets a gallery page showing its states in every theme, and a
check in step 3 keeps windows built from patterns. A screen that fits no
pattern is a proposal for a new pattern in kvit-ui, in the same way the
skill already treats a missing component. Risk: patterns designed too
early fit poorly, so they should be extracted from screens that already
work, and dots on the existing apps show which those are.

### 3.8 One rendering path

The approved picture and the shipped screen should be the same code.
kvit-ui's workflow today draws a proposed screen as HTML against generated
CSS, the owner approves the PNG, then an agent rebuilds it in Go. Colours
and sizes agree because CSS values are generated from Go design values, but
layout and behaviour are composed separately twice. That is the arrangement
that made web prototypes and Qt screens drift apart.

Three ways to remove the second build:

1. Mockups written in Go. A mockup is a function assembling real kvit-ui
   components with sample data, like a gallery specimen, rendered by the
   same headless path the gallery and `uitest` already use. Approved
   picture is then the code that ships. Exploring a radically different
   layout is slower than editing HTML. With page patterns, a screen is
   mostly slot assignments, keeping that cost small.
2. One declarative description with two renderers, one in HTML and one in
   Go. Moves drift from the agent into the renderers, which then have to be
   kept in step; a larger project.
3. Hosting a web interface inside the app. Suits Qt apps with QtWebEngine
   and does not apply to Go on unison.

The first fits the Kvit apps. HTML stays useful for one purpose: sketching
a new page pattern before it exists in kvit-ui, after which the pattern is
built in Go and judged in the gallery.

### 3.9 What is untested, and what it costs

- None of the method has been tried. Evidence-backed pieces: seeing the
  render, rubric of observed failures, positive examples, pairwise
  comparison, checks in the build rather than rules in prose, replays for
  completion and effort. Judgment-resting pieces: dots with confirmed
  reasons are enough feedback, the owner's share shrinks over rounds, page
  patterns fit the Kvit apps.
- Dots carry the owner's judgment and nobody else's. Suits apps with one
  expert user; would need other people's marks for wider-audience apps.
- Tokens spent mainly in writing scenarios (once per workflow) and in the
  review loop (once per change). Building page patterns and the markup mode
  are one-time engineering costs in kvit-ui.

### 3.10 First round on kvit-notes

kvit-notes is the app whose code has been read, so it is the natural first
test:

1. Draft `UX.md` for kvit-notes from the code and a short interview.
2. Turn three top workflows into scenarios that record completion, step
   counts and screenshots.
3. Show the owner the five screens the reviewer picks, for marking on
   screenshots, since the in-app markup mode does not exist yet.
4. Measure how many dots the agent's proposed reasons explained correctly,
   how many revisions the owner accepted, and how long the marking took.

That shows within a few days whether marking is as cheap as hoped and
whether the agent's proposed reasons are accurate enough, before kvit-ui
gets the markup mode or page patterns.

## 4. What has changed in kvit-ui since the review

Commit `bb47829` was the review's baseline. Since then, in this session's
settled commit `0f7d36c` (12 Nov 2026, +927 −49 across 6 files):

`sidebar.go` (+543 −40), `structure_test.go` (+235),
`agent/skills/kvit-ui/catalog.md`, `cmd/kvit-ui-gallery/catalog.go`,
`cmd/kvit-ui-gallery/gallery.go`,
`cmd/kvit-ui-gallery/specimens_structure.go` (+52).

What the sidebar now does:

- Group headings. `NewSidebarHeading(ui, "Folders")` goes among a
  sidebar's items and names the places under it. A screen reader hears it
  as a heading.
- Places inside places. `item.Add(children...)` puts folders inside a
  folder. The item then gets a small arrow before its symbol, each level
  indented one step.
- Opening and closing. Pointer: click on the arrow opens or closes without
  going there; click on the words goes there; double-click also opens or
  closes it. Keyboard: Right opens an item then steps into it; Left steps
  out to the parent then closes it. Expand and collapse both work through
  the screen reader.
- Selection inside a closed folder. The closed folder shows the accent
  bar, so the position stays visible.
- Narrow rail. Shows only the top level, without headings.
- Existing sidebars. A sidebar with no nested places keeps its layout and
  screen-reader tree.
- Clicked places take the keyboard. A clicked place takes keyboard focus
  even when it is not the selected one.
- Scrolling. The wheel scrolls the places with the same easing as other
  scrolling views. A scroll bar appears over the right-hand edge only while
  something is out of sight. The keyboard place is brought into view. Up,
  Down, Home and End still move through places and foot items as one group.
  The rail scrolls with wheel and keys but shows no bar. A sidebar can now
  be as short as its foot plus one row. Rows sit in an inner scrolling
  panel (`Sidebar.Rows()` returns them; `Position()` / `ScrollTo()` read
  and set the scroll).
- Gallery coverage. New specimens for grouped/nested sidebar and for a
  scrolling sidebar with a foot; KvitSidebar and KvitSidebarHeading pages;
  catalog regenerated. Screenshots in all four themes under `build/shots/`.

Tests added in `structure_test.go`:

- `TestASidebarGroupsItsPlacesAndOpensTheOnesThatHoldPlaces`
- `TestASidebarScrollsItsPlacesAboveItsFoot`

`./build.sh --test` passes with these tests, including the generated-file
checks (`catalog.md`, gallery catalog).

## 5. Still missing compared with kvit-notes' own sidebar

From `app/sidebar.go` in kvit-notes, not yet in kvit-ui:

- No folder colour mark field. `SetLead` replaces the symbol rather than
  adding a mark beside it.
- No drag of a folder onto another, no drop-target highlight.
- Counts stay pill `Badge`s; kvit-notes draws a plain number.
- No `Page Up` / `Page Down` in the sidebar.
- The scroll bar overlaps the right 10 px of rows while shown (about 2 px
  over a trailing count badge). Giving the bar its own strip would narrow
  every existing sidebar, so it was left overlapping.
- Code reading the sidebar's children directly for its rows would now find
  the inner scrolling panel instead.

Not checked: kvit-notes, kvit-cash and kvit-works were not compiled against
these two changes; nothing was run on a PC; pictures come from the headless
gallery on this machine.

kvit-notes uses kvit-ui from 12 November (`5faae058d552`, taken before the
sidebar work above), so it cannot use the new sidebar until that work reaches
kvit-ui's main branch. It still has no `.claude/skills` folder, still has the
letter buttons "B", "I", "U" in `app/toolbar.go:106` and
`editor/formatbar.go:60`, and its tests still call only `uitest.CheckRules`.
Bringing it onto kvit-ui with the sidebar, the skill and the checks remains the
first step of the trial in §3.10.

The sidebar work exists only on this track (`track/coding-agents-ux-8cfbb1`,
based on `bb47829`, 11 November). kvit-ui's main branch on GitHub has moved on
to 23 November (`ea7cba6`). `sidebar.go` has not changed there, but
`structure_test.go`, `agent/skills/kvit-ui/catalog.md`,
`cmd/kvit-ui-gallery/catalog.go` and `cmd/kvit-ui-gallery/specimens_structure.go`
have, so merging needs those four files reconciled.
