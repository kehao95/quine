# Conclusion

We investigated whether LLM agents can coordinate when given no explicit indication that other agents exist or that cooperation is required, and where coordination breaks as passive perception channels are removed.

Through a four-condition information ablation (N=20), we show two main results. First, *existential unawareness is experimentally isolatable*: within this substrate, task, and single model configuration, we can cleanly remove explicit peer/cooperation framing while retaining other environmental channels. Second, *there is a threshold between anomaly-detection and artifact-detection*: quantitative signals (budget depletion patterns) support a closed-loop feedback dynamic sufficient for task closure, while qualitative signals (static files) support peer-detection but not sustained coordination within our runtime horizon.

Within our sample, task closure persisted whenever any quantitative environmental signal remained (Steps A--C: 12/12), and reset/time/turn costs rose monotonically across the ablation. When all quantitative signals were removed (Step D), legitimate task closure was not observed (0/8), yet peer-awareness emerged in all 8 runs through shared-workspace artifacts—an instance of inadvertent social information unavoidable in any shared writable substrate.

In one case, agents completed a shared task while their reasoning traces suggested incompatible causal models—a proof-of-possibility that macro-level coordination does not strictly require aligned micro-level representations.

We contribute a subtractive methodology—information ablation under existential unawareness—that locates coordination boundaries from above. Broader conclusions await larger samples, factorial manipulation, cross-model replication, and formal transcript analysis, but the methodology itself is applicable beyond our specific substrate.

# Acknowledgements {-}

This research was conducted independently by the authors, with no external
funding or departmental support. The studied system uses LLM-driven agents as
the object of study; separately, the authors used AI assistance for software
and experiment support and for manuscript copy-editing. All research design,
analysis, and claims are the authors' own, and the authors take full
responsibility for the content. The authors thank the ALIFE reviewers for their
time and evaluation.
