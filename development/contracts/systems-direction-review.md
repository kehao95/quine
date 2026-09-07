# Systems-Direction Review

For a runtime or interface change, consider what it makes possible and how it
fits the operations around it. [Operating philosophy](../../Paper/philosophy/operating-philosophy.md)
and [development direction](../../Paper/philosophy/development-vector.md) offer
questions and existing design choices.

Useful review questions:

- Which process, state, channel, or access boundary changes?
- What existing mechanism can carry the operation?
- Which users or experiments depend on its current behavior?
- What does the new design clarify, enable, or simplify?
- What would show that the proposed behavior works?

A design may extend or reconsider the current architecture. Describe deliberate
contract changes and validate their consequences at the affected boundary.
