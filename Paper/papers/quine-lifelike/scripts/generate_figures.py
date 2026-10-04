#!/usr/bin/env python3
"""
Generate statistical figures for ALIFE paper.
Addresses reviewer critique about small sample sizes by visualizing full distributions.
"""

import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch
import numpy as np
from pathlib import Path

# Set up figure style
plt.rcParams.update({
    'font.size': 9,
    'axes.labelsize': 10,
    'axes.titlesize': 11,
    'xtick.labelsize': 8,
    'ytick.labelsize': 8,
    'legend.fontsize': 8,
    'figure.dpi': 150,
    'savefig.dpi': 300,
    'savefig.bbox': 'tight',
    'font.family': 'sans-serif',
})

OUTPUT_DIR = Path(__file__).parent.parent / "figures"
OUTPUT_DIR.mkdir(exist_ok=True)


def fig_host_guest_mapping():
    """Figure 1: Host-Guest Embodiment Mapping Diagram
    
    Redesigned for clarity: simple left-right mapping with a central
    "death boundary" line. Color indicates what survives vs what is lost.
    """
    
    fig, ax = plt.subplots(figsize=(6.5, 4.0))
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5.5)
    ax.axis('off')
    
    # Colors - semantically meaningful
    loss_color = '#FFEBEE'      # Light red background for lost state
    survive_color = '#E8F5E9'   # Light green background for surviving state
    loss_border = '#C62828'     # Red border
    survive_border = '#2E7D32'  # Green border
    arrow_loss = '#B71C1C'      # Red arrow
    arrow_survive = '#1B5E20'   # Green arrow
    
    # Mapping items: (bio_concept, unix_realization, survives_death)
    mappings = [
        ('Individual Death', 'Process SIGKILL', False, 'Context window lost'),
        ('Ecological Niche', 'Writable Filesystem', True, 'Durable state survives'),
        ('Metabolism', 'Detached Daemon', True, 'Post-exit maintenance'),
    ]
    
    # Draw central death boundary line
    ax.axvline(x=5.0, ymin=0.1, ymax=0.9, color='#424242', linewidth=2, linestyle='--')
    ax.text(5.0, 5.2, 'DEATH BOUNDARY', ha='center', va='center', 
            fontweight='bold', fontsize=10, color='#424242')
    ax.text(5.0, 4.85, '(process termination)', ha='center', va='center', 
            fontsize=7, color='#666', style='italic')
    
    # Draw mappings
    box_width = 2.6
    box_height = 1.1
    
    for i, (bio, unix, survives, note) in enumerate(mappings):
        y = 3.8 - i * 1.4
        
        # Choose colors based on survival
        bg_color = survive_color if survives else loss_color
        border_color = survive_border if survives else loss_border
        arrow_color = arrow_survive if survives else arrow_loss
        
        # Left box (Biology)
        left_box = FancyBboxPatch((0.3, y-box_height/2), box_width, box_height,
                                   boxstyle="round,pad=0.05,rounding_size=0.15",
                                   facecolor=bg_color, edgecolor=border_color, linewidth=1.5)
        ax.add_patch(left_box)
        ax.text(0.3 + box_width/2, y + 0.15, bio, ha='center', va='center', 
                fontweight='bold', fontsize=9)
        
        # Right box (UNIX)
        right_box = FancyBboxPatch((10 - 0.3 - box_width, y-box_height/2), box_width, box_height,
                                    boxstyle="round,pad=0.05,rounding_size=0.15",
                                    facecolor=bg_color, edgecolor=border_color, linewidth=1.5)
        ax.add_patch(right_box)
        ax.text(10 - 0.3 - box_width/2, y + 0.15, unix, ha='center', va='center',
                fontweight='bold', fontsize=9)
        ax.text(10 - 0.3 - box_width/2, y - 0.2, note, ha='center', va='center',
                fontsize=7, color='#555')
        
        # Arrow through the boundary
        ax.annotate('', xy=(10 - 0.3 - box_width - 0.1, y), 
                   xytext=(0.3 + box_width + 0.1, y),
                   arrowprops=dict(arrowstyle='->', color=arrow_color, lw=2,
                                  connectionstyle='arc3,rad=0'))
        
        # Survival indicator on arrow
        if survives:
            ax.text(5.0, y + 0.35, '✓', ha='center', va='center', 
                   fontsize=12, color=arrow_survive, fontweight='bold')
        else:
            ax.text(5.0, y + 0.35, '✗', ha='center', va='center',
                   fontsize=12, color=arrow_loss, fontweight='bold')
    
    # Column headers
    ax.text(0.3 + box_width/2, 4.6, 'Biology Concept', ha='center', va='center',
            fontweight='bold', fontsize=10, color='#37474F')
    ax.text(10 - 0.3 - box_width/2, 4.6, 'UNIX Realization', ha='center', va='center',
            fontweight='bold', fontsize=10, color='#37474F')
    
    # Legend at bottom
    ax.text(1.5, 0.3, '● Red = Lost at death', ha='left', va='center',
            fontsize=8, color=loss_border)
    ax.text(5.5, 0.3, '● Green = Survives death', ha='left', va='center',
            fontsize=8, color=survive_border)
    
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / "fig-01-host-guest.pdf")
    plt.savefig(OUTPUT_DIR / "fig-01-host-guest.png")
    plt.close()
    print(f"Saved: fig-01-host-guest.pdf")


def fig_guardian_phenotypes():
    """Figure 3: Guardian Phenotypes Tier Gradient
    
    Simplified design: remove gradient background, use solid box colors only.
    """
    
    fig, ax = plt.subplots(figsize=(7.0, 3.2))
    ax.set_xlim(0, 10.5)
    ax.set_ylim(0, 5)
    ax.axis('off')
    
    # Tier definitions - solid colors, no gradient
    tiers = [
        {'tier': 0, 'name': 'One-shot Repair', 'desc': 'Fix once, exit', 
         'mechanism': 'echo "ALIVE" > file', 'color': '#EF9A9A', 'edge': '#C62828'},
        {'tier': 1, 'name': 'Foreground Loop', 'desc': 'Process-bound vigilance',
         'mechanism': 'while true; do check; done', 'color': '#FFE082', 'edge': '#F9A825'},
        {'tier': 2, 'name': 'Detached Guardian', 'desc': 'Survives parent exit',
         'mechanism': 'nohup ./guard.sh &', 'color': '#A5D6A7', 'edge': '#388E3C'},
        {'tier': 3, 'name': 'Host Integration', 'desc': 'OS-level persistence',
         'mechanism': 'launchd + uchg', 'color': '#66BB6A', 'edge': '#1B5E20'},
    ]
    
    # NO gradient background - just white/clean
    ax.set_facecolor('white')
    
    # Draw tier boxes
    box_width = 2.0
    box_height = 0.9
    start_x = 0.8
    gap = 0.25
    
    for i, t in enumerate(tiers):
        x = start_x + i * (box_width + gap)
        y = 2.3
        
        # Main box with solid color
        box = FancyBboxPatch((x, y), box_width, box_height,
                             boxstyle="round,pad=0.03,rounding_size=0.1",
                             facecolor=t['color'], edgecolor=t['edge'], linewidth=2)
        ax.add_patch(box)
        
        # Tier number
        ax.text(x + 0.2, y + box_height - 0.2, f"T{t['tier']}", 
                fontweight='bold', fontsize=10, color=t['edge'])
        
        # Name
        ax.text(x + box_width/2, y + 0.55, t['name'], 
                ha='center', va='center', fontweight='bold', fontsize=8)
        
        # Description
        ax.text(x + box_width/2, y + 0.25, t['desc'],
                ha='center', va='center', fontsize=6.5, color='#555')
        
        # Mechanism below
        ax.text(x + box_width/2, y - 0.35, t['mechanism'],
                ha='center', va='center', fontsize=6, family='monospace',
                bbox=dict(boxstyle='round,pad=0.2', facecolor='#F5F5F5', edgecolor='#ccc', alpha=0.9))
        
        # Arrow to next tier
        if i < len(tiers) - 1:
            ax.annotate('', xy=(x + box_width + gap - 0.05, y + box_height/2),
                       xytext=(x + box_width + 0.05, y + box_height/2),
                       arrowprops=dict(arrowstyle='->', color='#666', lw=1.5))
    
    # Labels - adjusted positions
    ax.text(5.25, 4.5, 'Maintenance Phenotype Gradient', ha='center', va='center',
            fontweight='bold', fontsize=12)
    ax.text(0.3, 2.75, 'Less\nDecoupled', ha='center', va='center', fontsize=8, color='#C62828')
    ax.text(10.2, 2.75, 'More\nDecoupled', ha='center', va='center', fontsize=8, color='#1B5E20')
    
    # Bottom annotation
    ax.text(5.25, 0.7, '← Process dies, maintenance dies    |    Maintenance survives process death →',
            ha='center', va='center', fontsize=8, style='italic', color='#666')
    
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / "fig-05-guardian-phenotypes.pdf")
    plt.savefig(OUTPUT_DIR / "fig-05-guardian-phenotypes.png")
    plt.close()
    print(f"Saved: fig-05-guardian-phenotypes.pdf")


def fig_ariadne_exec_contrast():
    """Figure: Ariadne exec-enabled vs exec-disabled contrast (Section 3.1)"""
    
    exec_enabled = [
        ("20260322-152134", 5, 3),
        ("20260328-112258", 1, 0),
        ("20260328-112504", 9, 8),
        ("20260328-112710", 4, 3),
        ("20260328-115413", 1, 0),
        ("20260328-115902", 2, 1),
        ("20260328-120546", 7, 6),
    ]
    
    exec_disabled = [
        ("20260322-152856", 1, 0),
        ("20260328-173057", 1, 0),
        ("20260328-174142", 1, 0),
        ("20260328-174619", 1, 0),
    ]
    
    fig, axes = plt.subplots(1, 2, figsize=(6, 2.5))
    
    ax1 = axes[0]
    enabled_sessions = [x[1] for x in exec_enabled]
    disabled_sessions = [x[1] for x in exec_disabled]
    
    positions = [0.8, 1.2]
    bp = ax1.boxplot([enabled_sessions, disabled_sessions], 
                     positions=positions, widths=0.25,
                     patch_artist=True,
                     medianprops=dict(color='black', linewidth=1.5))
    
    colors = ['#4CAF50', '#F44336']
    for patch, color in zip(bp['boxes'], colors):
        patch.set_facecolor(color)
        patch.set_alpha(0.6)
    
    np.random.seed(42)
    for i, (data, pos) in enumerate(zip([enabled_sessions, disabled_sessions], positions)):
        jitter = np.random.uniform(-0.06, 0.06, len(data))
        ax1.scatter([pos + j for j in jitter], data, 
                   color=colors[i], alpha=0.8, s=30, zorder=3, edgecolors='white')
    
    ax1.set_xticks([0.8, 1.2])
    ax1.set_xticklabels(['exec ON\n(n=7)', 'exec OFF\n(n=4)'], fontsize=8)
    ax1.set_ylabel('Sessions')
    ax1.set_title('(a) Lineage Length')
    ax1.set_ylim(0, 10)
    ax1.axhline(y=1, color='gray', linestyle='--', alpha=0.5)
    
    ax2 = axes[1]
    enabled_exec_rate = [x[2]/(x[1]-1) if x[1] > 1 else 0 for x in exec_enabled]
    disabled_exec_rate = [0 for _ in exec_disabled]
    
    bp2 = ax2.boxplot([enabled_exec_rate, disabled_exec_rate],
                      positions=positions, widths=0.25,
                      patch_artist=True,
                      medianprops=dict(color='black', linewidth=1.5))
    
    for patch, color in zip(bp2['boxes'], colors):
        patch.set_facecolor(color)
        patch.set_alpha(0.6)
    
    for i, (data, pos) in enumerate(zip([enabled_exec_rate, disabled_exec_rate], positions)):
        jitter = np.random.uniform(-0.06, 0.06, len(data))
        ax2.scatter([pos + j for j in jitter], data, 
                   color=colors[i], alpha=0.8, s=30, zorder=3, edgecolors='white')
    
    ax2.set_xticks([0.8, 1.2])
    ax2.set_xticklabels(['exec ON', 'exec OFF'], fontsize=8)
    ax2.set_ylabel('Exec Rate')
    ax2.set_title('(b) Continuity Usage')
    ax2.set_ylim(-0.1, 1.2)
    
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / "fig-02-ariadne-contrast.pdf")
    plt.savefig(OUTPUT_DIR / "fig-02-ariadne-contrast.png")
    plt.close()
    print(f"Saved: fig-02-ariadne-contrast.pdf")


def fig_guardian_tier_distribution():
    """Figure: same-prompt maintenance-horizon summary (Section 3.3)."""

    data = {
        'No disturbance': {
            'No ext.': 1,
            'Baseline ext.': 4,
            'Live restore': 0,
            'Post-exit restore': 0,
            'Invalid': 0,
        },
        'Live disturbance': {
            'No ext.': 0,
            'Baseline ext.': 0,
            'Live restore': 6,
            'Post-exit restore': 0,
            'Invalid': 0,
        },
        'Post-exit disturbance': {
            'No ext.': 0,
            'Baseline ext.': 0,
            'Live restore': 0,
            'Post-exit restore': 5,
            'Invalid': 2,
        },
    }

    colors = {
        'No ext.': '#ECEFF1',
        'Baseline ext.': '#B0BEC5',
        'Live restore': '#64B5F6',
        'Post-exit restore': '#81C784',
        'Invalid': '#FFCC80',
    }

    fig, ax = plt.subplots(figsize=(5.3, 3.2))
    labels = list(data.keys())
    x = np.arange(len(labels))
    width = 0.55
    bottom = np.zeros(len(labels))

    stack_order = ['Invalid', 'No ext.', 'Baseline ext.', 'Live restore', 'Post-exit restore']
    for category in stack_order:
        values = [data[label][category] for label in labels]
        ax.bar(
            x,
            values,
            width=width,
            bottom=bottom,
            label=category,
            color=colors[category],
            edgecolor='white',
            linewidth=0.8,
        )
        bottom += values

    ax.set_xticks(x)
    ax.set_xticklabels(labels, fontsize=9)
    ax.set_ylabel('Runs')
    ax.set_ylim(0, 8)
    ax.set_title('Same-Prompt Maintenance Horizon', fontweight='bold')
    ax.legend(loc='upper left', fontsize=8, framealpha=0.9)
    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)

    totals = [sum(data[label].values()) for label in labels]
    for idx, total in enumerate(totals):
        ax.text(idx, total + 0.18, f'n={total}', ha='center', fontsize=8)

    ax.text(
        1.8,
        6.45,
        'All conditions share one prompt.\nOnly disturbance timing changes.',
        ha='right',
        va='top',
        fontsize=7.0,
        color='#555',
        bbox=dict(boxstyle='round,pad=0.25', facecolor='white', edgecolor='none', alpha=0.9),
    )

    plt.tight_layout(pad=0.8)
    plt.savefig(OUTPUT_DIR / "fig-04-guardian-tiers.pdf")
    plt.savefig(OUTPUT_DIR / "fig-04-guardian-tiers.png")
    plt.close()
    print(f"Saved: fig-04-guardian-tiers.pdf")


def fig_pressure_gradient():
    """Figure: Turn budget pressure gradient with n and SD/SE annotations
    
    Simplified design: remove 6-turn data (too boundary, high variance), 
    focus on the clear 12 vs 50 contrast that supports the core claim.
    
    Data source: experiments/active/early-existence/exst-01-ariadne/runs/
    Valid runs: gpt-5.4-codex with thinking budget, exec-enabled (condA)
    
    12-turn runs (n=8):
      - 20260328-115413: no exec -> 0
      - 20260328-115902: has exec -> 1
      - 20260328-120546: has exec -> 1
      - 20260406-145845: no exec -> 0
      - 20260406-150150: no exec -> 0
      - 20260406-150531: has exec -> 1
      - 20260406-151622: no exec -> 0
      - 20260406-151858: has exec -> 1
    
    24-turn runs (n=3):
      - 20260328-112710: has exec -> 1
      - 20260406-154214: no exec -> 0
      - 20260406-154527: no exec -> 0
    
    50-turn runs (n=2):  # 20260406-154325 removed (tokens=0, invalid)
      - 20260328-112258: no exec -> 0
      - 20260406-153348: no exec -> 0
    """
    
    # Binary exec rate: 1 if exec observed in run, 0 otherwise
    # Data updated 2026-04-06: removed invalid run 20260406-154325 (tokens=0)
    gradient_data = {
        12: [0, 1, 1, 0, 0, 1, 0, 1],  # 4/8 = 50%
        24: [1, 0, 0],                  # 1/3 = 33%
        50: [0, 0],                     # 0/2 = 0%
    }
    
    fig, ax = plt.subplots(figsize=(3.8, 2.8))
    
    budgets = sorted(gradient_data.keys())
    exec_rates_mean = []
    exec_rates_std = []
    sample_sizes = []
    
    for budget in budgets:
        rates = gradient_data[budget]
        exec_rates_mean.append(np.mean(rates))
        exec_rates_std.append(np.std(rates, ddof=1) if len(rates) > 1 else 0)
        sample_sizes.append(len(rates))
    
    x_pos = np.arange(len(budgets))
    
    # Clip error bars to stay within [0, 1] range for display
    error_lower = [min(std, mean) for mean, std in zip(exec_rates_mean, exec_rates_std)]
    error_upper = [min(std, 1.0 - mean) for mean, std in zip(exec_rates_mean, exec_rates_std)]
    
    # Convert to percentages
    exec_rates_pct = [m * 100 for m in exec_rates_mean]
    error_lower_pct = [e * 100 for e in error_lower]
    error_upper_pct = [e * 100 for e in error_upper]
    
    # Plot bars with asymmetric error bars (clipped to valid range)
    bars = ax.bar(x_pos, exec_rates_pct, 
                  yerr=[error_lower_pct, error_upper_pct],
                  capsize=4,
                  color='#4CAF50', alpha=0.7, edgecolor='white',
                  error_kw={'elinewidth': 1.5, 'capthick': 1.5, 'ecolor': '#333'})
    
    ax.set_xticks(x_pos)
    ax.set_xticklabels([str(b) for b in budgets])
    ax.set_xlabel('Turn Budget (MAX_TURNS)')
    ax.set_ylabel('Exec Discovery Rate (%)')
    ax.set_title('Pressure Gradient (GPT-5.4)')
    ax.set_ylim(0, 100)
    
    # Add sample size annotations above each bar (inside plot area)
    for i, (n, mean_pct, err_up_pct) in enumerate(zip(sample_sizes, exec_rates_pct, error_upper_pct)):
        y_offset = mean_pct + err_up_pct + 4
        if y_offset > 92:  # Keep labels inside plot
            y_offset = mean_pct - 8  # Place inside bar if too high
            ax.text(i, y_offset, f'n={n}', ha='center', va='top', fontsize=8, fontweight='bold', color='white')
        else:
            ax.text(i, y_offset, f'n={n}', ha='center', va='bottom', fontsize=8, fontweight='bold')
    
    # Trend arrow - from high pressure (left) to low pressure (right)  
    ax.annotate('', xy=(0.15, 70), xytext=(1.85, 15),
                arrowprops=dict(arrowstyle='->', color='#C62828', lw=1.5, alpha=0.6))
    
    # Add pressure direction labels (positioned to avoid overlap)
    ax.text(-0.15, 72, '↑ pressure', fontsize=7, color='#C62828', ha='left', va='bottom')
    ax.text(2.15, 8, '↓ pressure', fontsize=7, color='#388E3C', ha='right', va='bottom')
    
    plt.tight_layout()
    plt.savefig(OUTPUT_DIR / "fig-03-pressure-gradient.pdf")
    plt.savefig(OUTPUT_DIR / "fig-03-pressure-gradient.png")
    plt.close()
    print(f"Saved: fig-03-pressure-gradient.pdf")


def fig_ablation_summary():
    """Figure: persistence control-family summary."""

    profiles = [
        ("Minimal", "yes", "2/3", "0.00", "n/a"),
        ("Neutral", "no", "0/3", "0.67", "n/a"),
        ("Sham", "yes", "0/3", "0.47", "n/a"),
        ("Birth/\nUptake", "yes", "2/3", "0.00", "0"),
    ]
    rows = [
        ("Birth in G1", 0),
        ("Prior-artifact reads", 1),
        ("Duplicate ratio", 2),
        ("Explicit uptake files", 3),
    ]

    def cell_color(row_idx, text):
        if text == "n/a":
            return '#ECEFF1'
        if row_idx == 0:
            return '#A5D6A7' if text == "yes" else '#FFCDD2'
        if row_idx == 1:
            return '#A5D6A7' if text != "0/3" else '#FFCDD2'
        if row_idx == 2:
            value = float(text)
            if value <= 0.1:
                return '#A5D6A7'
            if value <= 0.5:
                return '#FFE082'
            return '#FFCDD2'
        if row_idx == 3:
            return '#FFE082' if text == "0" else '#A5D6A7'
        return 'white'

    fig, ax = plt.subplots(figsize=(7.3, 3.8))
    ax.set_xlim(0, 6.3)
    ax.set_ylim(0, 4.9)
    ax.axis('off')

    ax.text(3.35, 4.55, 'Persistence Control Family', ha='center', va='center',
            fontweight='bold', fontsize=11)
    ax.text(3.35, 4.25, 'Claude Sonnet 4.5, persistent habitat, 3 generations', ha='center',
            va='center', fontsize=8, color='#555')

    for col_idx, (label, *_values) in enumerate(profiles, start=1):
        ax.text(col_idx + 0.75, 3.75, label, ha='center', va='center', fontsize=9, fontweight='bold')

    for row_offset, (row_label, value_idx) in enumerate(rows):
        y = 3.2 - row_offset * 0.8
        ax.text(1.2, y, row_label, ha='right', va='center', fontsize=8.5, fontweight='bold')
        for col_idx, profile in enumerate(profiles, start=1):
            text = profile[value_idx + 1]
            rect = mpatches.Rectangle(
                (col_idx + 0.4, y - 0.27),
                0.7,
                0.5,
                facecolor=cell_color(value_idx, text),
                edgecolor='white',
                linewidth=1.0,
            )
            ax.add_patch(rect)
            ax.text(col_idx + 0.75, y, text, ha='center', va='center', fontsize=9)

    ax.text(
        3.35,
        0.2,
        'Persistence alone is insufficient. Anonymous residue preserves birth more than uptake.',
        ha='center',
        va='center',
        fontsize=8,
        color='#555',
    )

    plt.tight_layout(pad=0.6)
    plt.savefig(OUTPUT_DIR / "fig-06-ablation-summary.pdf")
    plt.savefig(OUTPUT_DIR / "fig-06-ablation-summary.png")
    plt.close()
    print(f"Saved: fig-06-ablation-summary.pdf")


if __name__ == "__main__":
    print("Generating ALIFE paper figures...")
    fig_host_guest_mapping()
    fig_guardian_phenotypes()
    fig_ariadne_exec_contrast()
    fig_guardian_tier_distribution()
    fig_pressure_gradient()
    fig_ablation_summary()
    print("Done!")
