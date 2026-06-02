<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GoQueue++ — Distributed Job Processing Engine</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;700&family=Syne:wght@400;600;700;800&display=swap" rel="stylesheet">
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --bg: #050a12;
    --bg2: #0c1220;
    --bg3: #111827;
    --border: #1e2d42;
    --border2: #243448;
    --blue: #3b82f6;
    --blue-dim: #1d4ed8;
    --blue-glow: rgba(59,130,246,0.15);
    --green: #22c55e;
    --green-dim: #15803d;
    --green-glow: rgba(34,197,94,0.12);
    --amber: #f59e0b;
    --amber-glow: rgba(245,158,11,0.12);
    --red: #ef4444;
    --purple: #8b5cf6;
    --text: #f1f5f9;
    --text2: #94a3b8;
    --text3: #475569;
    --mono: 'JetBrains Mono', monospace;
    --display: 'Syne', sans-serif;
  }

  html { scroll-behavior: smooth; }

  body {
    background: var(--bg);
    color: var(--text);
    font-family: var(--mono);
    font-size: 14px;
    line-height: 1.7;
    overflow-x: hidden;
  }

  /* ── ANIMATED GRID BACKGROUND ── */
  body::before {
    content: '';
    position: fixed;
    inset: 0;
    background-image:
      linear-gradient(rgba(59,130,246,0.03) 1px, transparent 1px),
      linear-gradient(90deg, rgba(59,130,246,0.03) 1px, transparent 1px);
    background-size: 40px 40px;
    z-index: 0;
    pointer-events: none;
  }

  /* ── SCANNING LINE ── */
  body::after {
    content: '';
    position: fixed;
    top: -100%;
    left: 0;
    right: 0;
    height: 200px;
    background: linear-gradient(transparent, rgba(59,130,246,0.04), transparent);
    animation: scan 8s linear infinite;
    z-index: 0;
    pointer-events: none;
  }

  @keyframes scan {
    0% { top: -200px; }
    100% { top: 100vh; }
  }

  .container {
    max-width: 900px;
    margin: 0 auto;
    padding: 0 24px 80px;
    position: relative;
    z-index: 1;
  }

  /* ── HERO ── */
  .hero {
    padding: 80px 0 60px;
    text-align: center;
    position: relative;
  }

  .hero-badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 6px 14px;
    border: 1px solid var(--border2);
    border-radius: 100px;
    font-size: 11px;
    color: var(--text2);
    margin-bottom: 28px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    background: var(--bg2);
    animation: fadeDown 0.6s ease both;
  }

  .hero-badge .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--green);
    box-shadow: 0 0 8px var(--green);
    animation: pulse 2s ease-in-out infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; box-shadow: 0 0 8px var(--green); }
    50% { opacity: 0.5; box-shadow: 0 0 4px var(--green); }
  }

  .hero h1 {
    font-family: var(--display);
    font-size: clamp(42px, 7vw, 72px);
    font-weight: 800;
    letter-spacing: -0.03em;
    line-height: 1;
    margin-bottom: 8px;
    animation: fadeUp 0.7s ease 0.1s both;
  }

  .hero h1 .brand { color: var(--blue); }
  .hero h1 .plus { color: var(--green); }

  .hero .version {
    font-family: var(--mono);
    font-size: 13px;
    color: var(--text3);
    margin-bottom: 20px;
    animation: fadeUp 0.7s ease 0.2s both;
  }

  .hero .tagline {
    font-family: var(--display);
    font-size: clamp(15px, 2.5vw, 18px);
    color: var(--text2);
    font-weight: 400;
    max-width: 520px;
    margin: 0 auto 36px;
    animation: fadeUp 0.7s ease 0.3s both;
  }

  .hero-stats {
    display: flex;
    justify-content: center;
    gap: 0;
    margin-bottom: 40px;
    border: 1px solid var(--border);
    border-radius: 12px;
    overflow: hidden;
    background: var(--bg2);
    animation: fadeUp 0.7s ease 0.4s both;
  }

  .hero-stat {
    flex: 1;
    padding: 20px 24px;
    border-right: 1px solid var(--border);
    text-align: center;
  }

  .hero-stat:last-child { border-right: none; }

  .hero-stat .num {
    font-family: var(--display);
    font-size: 28px;
    font-weight: 800;
    color: var(--blue);
    display: block;
    line-height: 1;
    margin-bottom: 4px;
  }

  .hero-stat .label {
    font-size: 11px;
    color: var(--text3);
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  @keyframes fadeUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @keyframes fadeDown {
    from { opacity: 0; transform: translateY(-12px); }
    to { opacity: 1; transform: translateY(0); }
  }

  /* ── SECTION ── */
  .section {
    margin-bottom: 56px;
    opacity: 0;
    transform: translateY(24px);
    transition: opacity 0.6s ease, transform 0.6s ease;
  }

  .section.visible {
    opacity: 1;
    transform: translateY(0);
  }

  .section-label {
    font-size: 10px;
    letter-spacing: 0.15em;
    text-transform: uppercase;
    color: var(--blue);
    margin-bottom: 6px;
    font-weight: 500;
  }

  .section-title {
    font-family: var(--display);
    font-size: 22px;
    font-weight: 700;
    color: var(--text);
    margin-bottom: 20px;
    letter-spacing: -0.01em;
  }

  /* ── INTRO CARD ── */
  .intro-card {
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 28px 32px;
    position: relative;
    overflow: hidden;
  }

  .intro-card::before {
    content: '';
    position: absolute;
    top: 0; left: 0; right: 0;
    height: 2px;
    background: linear-gradient(90deg, var(--blue), var(--green), var(--blue));
    background-size: 200% 100%;
    animation: shimmer 3s linear infinite;
  }

  @keyframes shimmer {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }

  .intro-card p {
    color: var(--text2);
    line-height: 1.8;
  }

  .intro-card p + p { margin-top: 12px; }

  .highlight { color: var(--blue); font-weight: 500; }
  .highlight-green { color: var(--green); font-weight: 500; }

  /* ── FEATURES GRID ── */
  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: 12px;
  }

  .feature-card {
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 18px 20px;
    display: flex;
    gap: 14px;
    align-items: flex-start;
    transition: border-color 0.25s, background 0.25s, transform 0.25s;
    cursor: default;
  }

  .feature-card:hover {
    border-color: var(--blue);
    background: var(--bg3);
    transform: translateY(-2px);
  }

  .feature-icon {
    font-size: 20px;
    flex-shrink: 0;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    background: rgba(59,130,246,0.08);
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid rgba(59,130,246,0.15);
  }

  .feature-name {
    font-family: var(--display);
    font-size: 13px;
    font-weight: 600;
    color: var(--text);
    margin-bottom: 4px;
  }

  .feature-desc {
    font-size: 12px;
    color: var(--text3);
    line-height: 1.5;
  }

  /* ── TECH STACK ── */
  .stack-grid {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 10px;
  }

  .stack-item {
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 16px 12px;
    text-align: center;
    transition: border-color 0.25s, transform 0.25s;
    cursor: default;
  }

  .stack-item:hover {
    border-color: var(--blue);
    transform: translateY(-3px);
  }

  .stack-item .icon { font-size: 22px; display: block; margin-bottom: 6px; }
  .stack-item .name { font-family: var(--display); font-size: 12px; font-weight: 600; color: var(--text); }
  .stack-item .role { font-size: 10px; color: var(--text3); margin-top: 2px; }

  /* ── ARCHITECTURE DIAGRAM ── */
  .arch-diagram {
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 32px 24px;
    overflow: hidden;
  }

  .arch-flow {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0;
  }

  .arch-node {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  .arch-box {
    padding: 12px 28px;
    border-radius: 8px;
    font-family: var(--display);
    font-size: 13px;
    font-weight: 600;
    letter-spacing: 0.01em;
    border: 1px solid;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    min-width: 180px;
    justify-content: center;
    position: relative;
    transition: box-shadow 0.3s;
  }

  .arch-box:hover { box-shadow: 0 0 20px rgba(59,130,246,0.2); }

  .box-client { background: rgba(30,41,63,0.8); border-color: var(--border2); color: var(--text2); }
  .box-api    { background: rgba(30,58,95,0.7); border-color: var(--blue); color: #93c5fd; }
  .box-pg     { background: rgba(28,42,58,0.7); border-color: #2563eb; color: #7dd3fc; }
  .box-redis  { background: rgba(32,13,13,0.7); border-color: var(--red); color: #fca5a5; }
  .box-worker { background: rgba(13,34,24,0.7); border-color: var(--green); color: #86efac; }
  .box-fail   { background: rgba(32,13,13,0.7); border-color: var(--red); color: #fca5a5; }
  .box-dlq    { background: rgba(28,18,38,0.7); border-color: var(--purple); color: #c4b5fd; }
  .box-reaper { background: rgba(28,18,38,0.7); border-color: var(--purple); color: #c4b5fd; }
  .box-obs    { background: rgba(13,28,53,0.7); border-color: var(--blue); color: #93c5fd; }

  .arch-arrow {
    width: 1px;
    height: 28px;
    background: var(--border2);
    position: relative;
    flex-shrink: 0;
  }

  .arch-arrow::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: -4px;
    width: 0;
    height: 0;
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-top: 6px solid var(--border2);
  }

  .arch-arrow.blue { background: var(--blue); }
  .arch-arrow.blue::after { border-top-color: var(--blue); }
  .arch-arrow.green { background: var(--green); }
  .arch-arrow.green::after { border-top-color: var(--green); }
  .arch-arrow.amber { background: var(--amber); }
  .arch-arrow.amber::after { border-top-color: var(--amber); }
  .arch-arrow.purple { background: var(--purple); opacity: 0.6; }
  .arch-arrow.purple::after { border-top-color: var(--purple); }

  /* Animated data packet on arrows */
  .arch-arrow .packet {
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--blue);
    animation: travel 2s ease-in-out infinite;
    box-shadow: 0 0 6px var(--blue);
  }

  .arch-arrow.green .packet { background: var(--green); box-shadow: 0 0 6px var(--green); animation-delay: 0.3s; }
  .arch-arrow.amber .packet { background: var(--amber); box-shadow: 0 0 6px var(--amber); animation-delay: 0.6s; }

  @keyframes travel {
    0% { top: 0; opacity: 0; }
    15% { opacity: 1; }
    85% { opacity: 1; }
    100% { top: 100%; opacity: 0; }
  }

  /* Queue fan-out row */
  .queue-row {
    display: flex;
    gap: 8px;
    justify-content: center;
    width: 100%;
    margin: 4px 0;
  }

  .queue-pill {
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.04em;
    border: 1px solid;
    font-family: var(--display);
  }

  .q-high   { background: rgba(120,53,15,0.3); border-color: var(--amber); color: #fcd34d; }
  .q-medium { background: rgba(120,53,15,0.2); border-color: #b45309; color: #fde68a; }
  .q-low    { background: rgba(120,53,15,0.1); border-color: #78350f; color: #fef3c7; }

  .fan-lines {
    display: flex;
    justify-content: center;
    width: 100%;
    height: 24px;
    position: relative;
  }

  .fan-lines svg {
    position: absolute;
    top: 0; left: 0; right: 0; bottom: 0;
    width: 100%; height: 100%;
    overflow: visible;
  }

  /* Worker pool */
  .worker-pool {
    border: 1px solid var(--green);
    border-radius: 10px;
    background: rgba(13,34,24,0.5);
    padding: 14px 20px;
    width: 100%;
    max-width: 400px;
  }

  .pool-label {
    font-size: 10px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--green);
    text-align: center;
    margin-bottom: 10px;
    font-weight: 600;
  }

  .worker-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .worker-item {
    background: rgba(20,83,45,0.4);
    border: 1px solid rgba(34,197,94,0.2);
    border-radius: 6px;
    padding: 7px 14px;
    font-size: 12px;
    color: #86efac;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .worker-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--green);
    box-shadow: 0 0 6px var(--green);
    flex-shrink: 0;
  }

  .worker-item:nth-child(1) .worker-dot { animation: blink 1.8s ease-in-out infinite; }
  .worker-item:nth-child(2) .worker-dot { animation: blink 1.8s ease-in-out 0.6s infinite; }
  .worker-item:nth-child(3) .worker-dot { animation: blink 1.8s ease-in-out 1.2s infinite; }

  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  .worker-bar {
    flex: 1;
    height: 4px;
    background: rgba(34,197,94,0.1);
    border-radius: 2px;
    overflow: hidden;
  }

  .worker-bar-fill {
    height: 100%;
    background: var(--green);
    border-radius: 2px;
    animation: workload 3s ease-in-out infinite;
  }

  .worker-item:nth-child(2) .worker-bar-fill { animation-delay: 1s; }
  .worker-item:nth-child(3) .worker-bar-fill { animation-delay: 2s; }

  @keyframes workload {
    0% { width: 20%; }
    50% { width: 85%; }
    100% { width: 35%; }
  }

  /* Failure row */
  .failure-row {
    display: flex;
    gap: 8px;
    justify-content: center;
    flex-wrap: wrap;
    width: 100%;
    max-width: 520px;
  }

  .fail-box {
    padding: 10px 14px;
    border-radius: 8px;
    font-size: 11px;
    border: 1px solid;
    font-family: var(--display);
    font-weight: 600;
    text-align: center;
    flex: 1;
    min-width: 120px;
  }

  .fb-retry  { background: rgba(120,53,15,0.2); border-color: var(--amber); color: #fcd34d; }
  .fb-dlq    { background: rgba(32,13,13,0.3); border-color: var(--red); color: #fca5a5; }
  .fb-vis    { background: rgba(13,28,53,0.3); border-color: var(--blue); color: #93c5fd; }

  .fail-box .sub {
    font-family: var(--mono);
    font-size: 10px;
    font-weight: 400;
    color: var(--text3);
    margin-top: 3px;
    display: block;
  }

  /* Recovery arrow */
  .recovery-note {
    font-size: 11px;
    color: var(--purple);
    text-align: center;
    padding: 6px 16px;
    border: 1px dashed rgba(139,92,246,0.3);
    border-radius: 20px;
    background: rgba(139,92,246,0.05);
  }

  /* ── CODE BLOCK ── */
  .code-block {
    background: #050a12;
    border: 1px solid var(--border);
    border-radius: 10px;
    overflow: hidden;
  }

  .code-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 16px;
    background: var(--bg2);
    border-bottom: 1px solid var(--border);
  }

  .code-dots { display: flex; gap: 6px; }
  .code-dot {
    width: 10px; height: 10px; border-radius: 50%;
  }
  .cd-red { background: #ef4444; }
  .cd-amber { background: #f59e0b; }
  .cd-green { background: #22c55e; }

  .code-filename {
    font-size: 11px;
    color: var(--text3);
    letter-spacing: 0.05em;
  }

  .copy-btn {
    background: none;
    border: 1px solid var(--border2);
    color: var(--text3);
    font-size: 11px;
    font-family: var(--mono);
    padding: 4px 10px;
    border-radius: 5px;
    cursor: pointer;
    transition: border-color 0.2s, color 0.2s;
  }

  .copy-btn:hover { border-color: var(--blue); color: var(--blue); }
  .copy-btn.copied { border-color: var(--green); color: var(--green); }

  .code-body {
    padding: 20px 20px;
    font-size: 13px;
    line-height: 1.8;
    color: #8892a4;
    overflow-x: auto;
    white-space: pre;
  }

  .code-body .kw { color: #c792ea; }
  .code-body .fn { color: #82aaff; }
  .code-body .str { color: #c3e88d; }
  .code-body .num { color: #f78c6c; }
  .code-body .cm { color: #546e7a; font-style: italic; }
  .code-body .key { color: #89ddff; }
  .code-body .val { color: #f8fafc; }

  /* ── LOG DEMO ── */
  .log-demo {
    background: #050a12;
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 20px;
    font-size: 12px;
    overflow: hidden;
  }

  .log-line {
    display: flex;
    gap: 12px;
    padding: 4px 0;
    opacity: 0;
    animation: logAppear 0.4s ease forwards;
    border-bottom: 1px solid rgba(30,45,66,0.5);
  }

  .log-line:last-child { border-bottom: none; }

  @keyframes logAppear {
    from { opacity: 0; transform: translateX(-8px); }
    to { opacity: 1; transform: translateX(0); }
  }

  .log-time { color: var(--text3); flex-shrink: 0; }
  .log-level-info { color: var(--blue); font-weight: 600; flex-shrink: 0; width: 38px; }
  .log-level-warn { color: var(--amber); font-weight: 600; flex-shrink: 0; width: 38px; }
  .log-level-err { color: var(--red); font-weight: 600; flex-shrink: 0; width: 38px; }
  .log-corr { color: var(--amber); }
  .log-event { color: var(--green); }
  .log-rest { color: var(--text3); }

  /* ── STATS ROW ── */
  .stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
  }

  .stat-card {
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 20px 16px;
    text-align: center;
    transition: border-color 0.25s;
  }

  .stat-card:hover { border-color: var(--blue); }

  .stat-num {
    font-family: var(--display);
    font-size: 30px;
    font-weight: 800;
    line-height: 1;
    margin-bottom: 6px;
    display: block;
  }

  .stat-label { font-size: 11px; color: var(--text3); text-transform: uppercase; letter-spacing: 0.08em; }

  /* ── QUICK START ── */
  .steps { display: flex; flex-direction: column; gap: 14px; }

  .step {
    display: flex;
    gap: 16px;
    align-items: flex-start;
    background: var(--bg2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 18px 20px;
    transition: border-color 0.25s;
  }

  .step:hover { border-color: var(--blue); }

  .step-num {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: rgba(59,130,246,0.15);
    border: 1px solid var(--blue);
    color: var(--blue);
    font-family: var(--display);
    font-size: 13px;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .step-content { flex: 1; }
  .step-title { font-family: var(--display); font-size: 14px; font-weight: 600; color: var(--text); margin-bottom: 8px; }

  /* ── FILE TREE ── */
  .file-tree {
    background: #050a12;
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 20px 24px;
    font-size: 13px;
    line-height: 2;
  }

  .tree-root { color: var(--blue); font-weight: 600; }
  .tree-dir { color: var(--text2); }
  .tree-file { color: var(--text3); }
  .tree-comment { color: #334155; font-style: italic; }
  .tree-indent { padding-left: 20px; display: block; }
  .tree-indent2 { padding-left: 40px; display: block; }

  /* ── SERVICES TABLE ── */
  .services-table {
    width: 100%;
    border-collapse: collapse;
  }

  .services-table th {
    text-align: left;
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text3);
    padding: 8px 16px;
    border-bottom: 1px solid var(--border);
    font-weight: 500;
  }

  .services-table td {
    padding: 12px 16px;
    border-bottom: 1px solid rgba(30,45,66,0.5);
    font-size: 13px;
    color: var(--text2);
  }

  .services-table tr:last-child td { border-bottom: none; }
  .services-table tr:hover td { background: rgba(59,130,246,0.04); color: var(--text); }

  .service-name { color: var(--text); font-weight: 500; }
  .service-url { color: var(--blue); font-family: var(--mono); font-size: 12px; }

  /* ── FOOTER ── */
  .footer {
    margin-top: 80px;
    padding: 40px 0 20px;
    border-top: 1px solid var(--border);
    text-align: center;
  }

  .footer-brand {
    font-family: var(--display);
    font-size: 20px;
    font-weight: 800;
    color: var(--blue);
    margin-bottom: 8px;
  }

  .footer-sub { font-size: 12px; color: var(--text3); }

  .footer-chips {
    display: flex;
    gap: 8px;
    justify-content: center;
    flex-wrap: wrap;
    margin-top: 16px;
  }

  .chip {
    padding: 4px 12px;
    border-radius: 100px;
    font-size: 11px;
    border: 1px solid var(--border2);
    color: var(--text3);
    background: var(--bg2);
    font-family: var(--mono);
  }

  /* ── SCROLLBAR ── */
  ::-webkit-scrollbar { width: 6px; height: 6px; }
  ::-webkit-scrollbar-track { background: var(--bg); }
  ::-webkit-scrollbar-thumb { background: var(--border2); border-radius: 3px; }
  ::-webkit-scrollbar-thumb:hover { background: var(--blue-dim); }

  /* ── RESPONSIVE ── */
  @media (max-width: 640px) {
    .stack-grid { grid-template-columns: repeat(3, 1fr); }
    .stats-row { grid-template-columns: repeat(2, 1fr); }
    .hero-stats { flex-direction: column; }
    .hero-stat { border-right: none; border-bottom: 1px solid var(--border); }
    .hero-stat:last-child { border-bottom: none; }
  }
</style>
</head>
<body>

<div class="container">

  <!-- ═══ HERO ═══ -->
  <div class="hero">
    <div class="hero-badge">
      <span class="dot"></span>
      v1.0 — Production Ready
    </div>
    <h1><span class="brand">GoQueue</span><span class="plus">++</span></h1>
    <p class="version">Built in Public · 12 Posts Later</p>
    <p class="tagline">A fault-tolerant, distributed job processing engine built in Go. Built for teams who can't afford to lose a single job.</p>
    <div class="hero-stats">
      <div class="hero-stat">
        <span class="num" data-count="10000">0</span>
        <span class="label">Virtual Users</span>
      </div>
      <div class="hero-stat">
        <span class="num" data-count="40000">0</span>
        <span class="label">Requests Tested</span>
      </div>
      <div class="hero-stat">
        <span class="num">3</span>
        <span class="label">Priority Queues</span>
      </div>
      <div class="hero-stat">
        <span class="num">100%</span>
        <span class="label">Job Recovery</span>
      </div>
    </div>
  </div>

  <!-- ═══ WHAT IS GOQUEUE ═══ -->
  <div class="section">
    <p class="section-label">Overview</p>
    <h2 class="section-title">What is GoQueue?</h2>
    <div class="intro-card">
      <p>Processing everything inside a synchronous HTTP request is fragile — timeouts kill long-running work, retries duplicate side effects, and a single crash can silently swallow jobs.</p>
      <p>GoQueue solves this by <span class="highlight">decoupling job creation from job execution</span>. Every job is persisted to <span class="highlight">PostgreSQL</span> the moment it's created, pushed to <span class="highlight">Redis</span> for fast async dispatch, and executed by a pool of concurrent <span class="highlight-green">Go workers</span> — with full crash recovery, priority scheduling, and end-to-end tracing built in.</p>
    </div>
  </div>

  <!-- ═══ FEATURES ═══ -->
  <div class="section">
    <p class="section-label">Features</p>
    <h2 class="section-title">Core Capabilities</h2>
    <div class="features-grid">
      <div class="feature-card">
        <div class="feature-icon">🗄️</div>
        <div>
          <div class="feature-name">Dual-Layer Persistence</div>
          <div class="feature-desc">PostgreSQL for durability, Redis for throughput. Jobs survive restarts.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">👷</div>
        <div>
          <div class="feature-name">Concurrent Worker Pool</div>
          <div class="feature-desc">Scalable goroutine-based workers processing jobs in parallel.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🔄</div>
        <div>
          <div class="feature-name">Smart Retry + Backoff</div>
          <div class="feature-desc">Non-blocking delayed queues with exponential backoff to protect downstream systems.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">☠️</div>
        <div>
          <div class="feature-name">Dead Letter Queue</div>
          <div class="feature-desc">Permanently failed jobs routed to DLQ for inspection and replay.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">💀</div>
        <div>
          <div class="feature-name">Crash Recovery</div>
          <div class="feature-desc">Visibility timeouts + Reaper service recover jobs from crashed workers automatically.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🎯</div>
        <div>
          <div class="feature-name">Priority Dispatching</div>
          <div class="feature-desc">Three-tier queue dispatching: High, Medium, and Low priority execution.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🔑</div>
        <div>
          <div class="feature-name">Idempotency Keys</div>
          <div class="feature-desc">Safe job creation — retrying the same request never creates duplicates.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🔍</div>
        <div>
          <div class="feature-name">Correlation IDs</div>
          <div class="feature-desc">Every job is traceable end-to-end via structured JSON logs.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🛑</div>
        <div>
          <div class="feature-name">Backpressure Load Shedding</div>
          <div class="feature-desc">Returns HTTP 429 when queue capacity is exceeded to protect system stability.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">📊</div>
        <div>
          <div class="feature-name">Full Observability</div>
          <div class="feature-desc">Prometheus metrics + Grafana dashboards out of the box.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🧹</div>
        <div>
          <div class="feature-name">Graceful Shutdown</div>
          <div class="feature-desc">Workers drain in-flight jobs safely on SIGTERM/SIGINT signals.</div>
        </div>
      </div>
      <div class="feature-card">
        <div class="feature-icon">⚡</div>
        <div>
          <div class="feature-name">High Throughput</div>
          <div class="feature-desc">Validated at 40,000+ requests with 10,000 virtual users under k6 load testing.</div>
        </div>
      </div>
    </div>
  </div>

  <!-- ═══ TECH STACK ═══ -->
  <div class="section">
    <p class="section-label">Stack</p>
    <h2 class="section-title">Tech Stack</h2>
    <div class="stack-grid">
      <div class="stack-item">
        <span class="icon">🐹</span>
        <div class="name">Go</div>
        <div class="role">Language</div>
      </div>
      <div class="stack-item">
        <span class="icon">🐘</span>
        <div class="name">PostgreSQL</div>
        <div class="role">Persistence</div>
      </div>
      <div class="stack-item">
        <span class="icon">🔴</span>
        <div class="name">Redis</div>
        <div class="role">Queue / Cache</div>
      </div>
      <div class="stack-item">
        <span class="icon">📈</span>
        <div class="name">Prometheus</div>
        <div class="role">Metrics</div>
      </div>
      <div class="stack-item">
        <span class="icon">📊</span>
        <div class="name">Grafana</div>
        <div class="role">Dashboards</div>
      </div>
    </div>
  </div>

  <!-- ═══ ARCHITECTURE ═══ -->
  <div class="section">
    <p class="section-label">Design</p>
    <h2 class="section-title">Architecture Overview</h2>
    <div class="arch-diagram">
      <div class="arch-flow">

        <div class="arch-node">
          <div class="arch-box box-client">Client Request</div>
        </div>
        <div class="arch-arrow blue"><div class="packet"></div></div>

        <div class="arch-node">
          <div class="arch-box box-api">API Service</div>
        </div>
        <div class="arch-arrow blue"><div class="packet"></div></div>

        <div class="arch-node">
          <div class="arch-box box-pg">PostgreSQL <span style="font-size:10px;color:#475569;font-weight:400">← persists job</span></div>
        </div>
        <div class="arch-arrow blue"><div class="packet"></div></div>

        <div class="arch-node">
          <div class="arch-box box-redis">Redis</div>
        </div>

        <div style="height:16px;"></div>

        <div class="queue-row">
          <div class="queue-pill q-high">⬆ High Priority</div>
          <div class="queue-pill q-medium">→ Medium Priority</div>
          <div class="queue-pill q-low">⬇ Low Priority</div>
        </div>

        <div class="arch-arrow green" style="margin-top:4px;"><div class="packet"></div></div>

        <div class="arch-node" style="width:100%;max-width:400px;">
          <div class="worker-pool">
            <div class="pool-label">Worker Pool</div>
            <div class="worker-list">
              <div class="worker-item">
                <div class="worker-dot"></div>
                Worker 1
                <div class="worker-bar"><div class="worker-bar-fill"></div></div>
                <span style="font-size:10px;color:#64748b;">active</span>
              </div>
              <div class="worker-item">
                <div class="worker-dot"></div>
                Worker 2
                <div class="worker-bar"><div class="worker-bar-fill"></div></div>
                <span style="font-size:10px;color:#64748b;">active</span>
              </div>
              <div class="worker-item">
                <div class="worker-dot"></div>
                Worker 3
                <div class="worker-bar"><div class="worker-bar-fill"></div></div>
                <span style="font-size:10px;color:#64748b;">active</span>
              </div>
            </div>
          </div>
        </div>

        <div class="arch-arrow amber"><div class="packet"></div></div>

        <div class="failure-row">
          <div class="fail-box fb-retry">
            Retry Queue
            <span class="sub">Exponential Backoff</span>
          </div>
          <div class="fail-box fb-dlq">
            Dead Letter Queue
            <span class="sub">Failed Jobs</span>
          </div>
          <div class="fail-box fb-vis">
            Visibility Timeout
            <span class="sub">Job Lease Protection</span>
          </div>
        </div>

        <div class="arch-arrow purple"><div class="packet"></div></div>

        <div class="arch-node" style="width:100%;max-width:280px;">
          <div class="arch-box box-reaper" style="width:100%;justify-content:center;">
            Reaper Service <span style="font-size:10px;color:#6d28d9;margin-left:6px;">crash recovery</span>
          </div>
        </div>

        <div style="height:12px;"></div>
        <div class="recovery-note">↻ crashed worker jobs are returned to the queue automatically</div>
        <div style="height:16px;"></div>

        <div class="arch-node" style="width:100%;max-width:360px;">
          <div class="arch-box box-obs" style="width:100%;justify-content:center;">
            Prometheus · Grafana · Correlation IDs
          </div>
        </div>

      </div>
    </div>
  </div>

  <!-- ═══ STRUCTURED LOGS ═══ -->
  <div class="section">
    <p class="section-label">Observability</p>
    <h2 class="section-title">Structured Logs</h2>
    <div class="log-demo" id="logDemo">
      <div class="log-line" style="animation-delay:0.1s">
        <span class="log-time">10:24:01.142</span>
        <span class="log-level-info">INFO</span>
        <span class="log-corr">REQ-8472</span>
        <span class="log-event">job_created</span>
        <span class="log-rest">queue=high job_id=a1b2c3</span>
      </div>
      <div class="log-line" style="animation-delay:0.5s">
        <span class="log-time">10:24:01.158</span>
        <span class="log-level-info">INFO</span>
        <span class="log-corr">REQ-8472</span>
        <span class="log-event">job_enqueued</span>
        <span class="log-rest">redis_key=goqueue:high:a1b2c3</span>
      </div>
      <div class="log-line" style="animation-delay:0.9s">
        <span class="log-time">10:24:01.204</span>
        <span class="log-level-info">INFO</span>
        <span class="log-corr">REQ-8472</span>
        <span class="log-event">job_dequeued</span>
        <span class="log-rest">worker=worker-1 attempt=1</span>
      </div>
      <div class="log-line" style="animation-delay:1.3s">
        <span class="log-time">10:24:01.347</span>
        <span class="log-level-info">INFO</span>
        <span class="log-corr">REQ-8472</span>
        <span class="log-event">job_processed</span>
        <span class="log-rest">duration_ms=143 status=success</span>
      </div>
      <div class="log-line" style="animation-delay:1.7s">
        <span class="log-time">10:24:03.021</span>
        <span class="log-level-warn">WARN</span>
        <span class="log-corr">REQ-9103</span>
        <span class="log-event">job_failed</span>
        <span class="log-rest">attempt=2 next_retry_in=4s</span>
      </div>
      <div class="log-line" style="animation-delay:2.1s">
        <span class="log-time">10:24:07.003</span>
        <span class="log-level-err">ERR</span>
        <span class="log-corr">REQ-9103</span>
        <span class="log-event">job_dlq</span>
        <span class="log-rest">max_attempts=3 moved_to=dlq</span>
      </div>
    </div>
  </div>

  <!-- ═══ QUICK START ═══ -->
  <div class="section">
    <p class="section-label">Setup</p>
    <h2 class="section-title">Quick Start</h2>
    <div class="steps">
      <div class="step">
        <div class="step-num">1</div>
        <div class="step-content">
          <div class="step-title">Clone the repository</div>
          <div class="code-block">
            <div class="code-header">
              <div class="code-dots"><div class="code-dot cd-red"></div><div class="code-dot cd-amber"></div><div class="code-dot cd-green"></div></div>
              <span class="code-filename">terminal</span>
              <button class="copy-btn" onclick="copyCode(this, 'git clone https://github.com/bhusaremayur/goqueue.git\ncd goqueue')">copy</button>
            </div>
            <div class="code-body">git clone https://github.com/bhusaremayur/goqueue.git
cd goqueue</div>
          </div>
        </div>
      </div>

      <div class="step">
        <div class="step-num">2</div>
        <div class="step-content">
          <div class="step-title">Start all services via Docker Compose</div>
          <div class="code-block">
            <div class="code-header">
              <div class="code-dots"><div class="code-dot cd-red"></div><div class="code-dot cd-amber"></div><div class="code-dot cd-green"></div></div>
              <span class="code-filename">terminal</span>
              <button class="copy-btn" onclick="copyCode(this, 'docker-compose -f deployments/docker-compose.yml up -d\n\n# Or with Make:\nmake up')">copy</button>
            </div>
            <div class="code-body"><span class="cm"># Spins up API, workers, PostgreSQL, Redis, Prometheus, Grafana</span>
docker-compose -f deployments/docker-compose.yml up -d

<span class="cm"># Or with Make:</span>
make up</div>
          </div>
        </div>
      </div>

      <div class="step">
        <div class="step-num">3</div>
        <div class="step-content">
          <div class="step-title">Verify services are running</div>
          <table class="services-table" style="margin-top:4px;">
            <thead>
              <tr>
                <th>Service</th>
                <th>URL</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td class="service-name">Go API Server</td>
                <td class="service-url">http://localhost:8080</td>
              </tr>
              <tr>
                <td class="service-name">Grafana Dashboard</td>
                <td class="service-url">http://localhost:3000</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <!-- ═══ REPO STRUCTURE ═══ -->
  <div class="section">
    <p class="section-label">Codebase</p>
    <h2 class="section-title">Repository Structure</h2>
    <div class="code-block">
      <div class="code-header">
        <div class="code-dots"><div class="code-dot cd-red"></div><div class="code-dot cd-amber"></div><div class="code-dot cd-green"></div></div>
        <span class="code-filename">goqueue/</span>
        <button class="copy-btn" onclick="copyCode(this, 'goqueue/\n├── cmd/\n│   ├── api/            Entry point for the API server\n│   ├── worker/         Entry point for the worker process\n│   └── example-app/    Example usage\n│\n├── internal/\n│   ├── job/            Core job domain logic\n│   ├── storage/\n│   │   ├── postgres/   PostgreSQL repositories\n│   │   │   └── migrations/  SQL schema migrations\n│   │   └── redis/      Redis queue implementation\n│   ├── scheduler/      Job scheduling logic\n│   └── reaper/         Crash recovery service\n│\n├── pkg/\n│   ├── logging/        Structured JSON logger\n│   ├── metrics/        Prometheus metric definitions\n│   ├── middleware/      HTTP middleware (correlation IDs, auth)\n│   └── retry/          Exponential backoff policies\n│\n└── deployments/\n    ├── docker-compose.yml\n    └── prometheus.yml')">copy</button>
      </div>
      <div class="file-tree">
<span class="tree-root">goqueue/</span>
<span class="tree-indent"><span class="tree-dir">├── cmd/</span></span>
<span class="tree-indent2"><span class="tree-dir">├── api/</span>           <span class="tree-comment">Entry point for the API server</span></span>
<span class="tree-indent2"><span class="tree-dir">├── worker/</span>        <span class="tree-comment">Entry point for the worker process</span></span>
<span class="tree-indent2"><span class="tree-dir">└── example-app/</span>   <span class="tree-comment">Example usage</span></span>
<span class="tree-indent" style="height:8px;display:block;"></span>
<span class="tree-indent"><span class="tree-dir">├── internal/</span></span>
<span class="tree-indent2"><span class="tree-dir">├── job/</span>           <span class="tree-comment">Core job domain logic</span></span>
<span class="tree-indent2"><span class="tree-dir">├── storage/</span></span>
<span class="tree-indent2" style="padding-left:60px;"><span class="tree-dir">├── postgres/</span>  <span class="tree-comment">PostgreSQL repositories</span></span>
<span class="tree-indent2" style="padding-left:80px;"><span class="tree-dir">└── migrations/</span> <span class="tree-comment">SQL schema migrations</span></span>
<span class="tree-indent2" style="padding-left:60px;"><span class="tree-dir">└── redis/</span>     <span class="tree-comment">Redis queue implementation</span></span>
<span class="tree-indent2"><span class="tree-dir">├── scheduler/</span>     <span class="tree-comment">Job scheduling logic</span></span>
<span class="tree-indent2"><span class="tree-dir">└── reaper/</span>        <span class="tree-comment">Crash recovery service</span></span>
<span class="tree-indent" style="height:8px;display:block;"></span>
<span class="tree-indent"><span class="tree-dir">├── pkg/</span></span>
<span class="tree-indent2"><span class="tree-dir">├── logging/</span>       <span class="tree-comment">Structured JSON logger</span></span>
<span class="tree-indent2"><span class="tree-dir">├── metrics/</span>       <span class="tree-comment">Prometheus metric definitions</span></span>
<span class="tree-indent2"><span class="tree-dir">├── middleware/</span>    <span class="tree-comment">HTTP middleware (correlation IDs, auth)</span></span>
<span class="tree-indent2"><span class="tree-dir">└── retry/</span>         <span class="tree-comment">Exponential backoff policies</span></span>
<span class="tree-indent" style="height:8px;display:block;"></span>
<span class="tree-indent"><span class="tree-dir">└── deployments/</span></span>
<span class="tree-indent2"><span class="tree-file">├── docker-compose.yml</span></span>
<span class="tree-indent2"><span class="tree-file">└── prometheus.yml</span></span>
      </div>
    </div>
  </div>

  <!-- ═══ LOAD TESTING ═══ -->
  <div class="section">
    <p class="section-label">Performance</p>
    <h2 class="section-title">Load Testing Results</h2>
    <div class="stats-row" style="margin-bottom:16px;">
      <div class="stat-card">
        <span class="stat-num" style="color:var(--blue);" data-count="10000">0</span>
        <span class="stat-label">Virtual Users</span>
      </div>
      <div class="stat-card">
        <span class="stat-num" style="color:var(--amber);" data-count="40000">0</span>
        <span class="stat-label">Requests Tested</span>
      </div>
      <div class="stat-card">
        <span class="stat-num" style="color:var(--green);">✓</span>
        <span class="stat-label">Backpressure</span>
      </div>
      <div class="stat-card">
        <span class="stat-num" style="color:var(--green);">✓</span>
        <span class="stat-label">DLQ Routing</span>
      </div>
    </div>
    <div class="code-block">
      <div class="code-header">
        <div class="code-dots"><div class="code-dot cd-red"></div><div class="code-dot cd-amber"></div><div class="code-dot cd-green"></div></div>
        <span class="code-filename">k6 results</span>
      </div>
      <div class="code-body"><span class="key">virtual_users</span>    <span class="val">10,000</span>
<span class="key">requests_tested</span>  <span class="val">40,000+</span>
<span class="key">backpressure</span>     <span class="val">✓ HTTP 429 fired correctly at queue capacity</span>
<span class="key">crash_recovery</span>   <span class="val">✓ Jobs recovered after worker kill -9</span>
<span class="key">dlq_routing</span>      <span class="val">✓ Permanently failed jobs isolated cleanly</span>
<span class="key">idempotency</span>      <span class="val">✓ Duplicate requests created zero duplicate jobs</span></div>
    </div>
  </div>

  <!-- ═══ CONTRIBUTING ═══ -->
  <div class="section">
    <p class="section-label">Community</p>
    <h2 class="section-title">Contributing</h2>
    <div class="intro-card">
      <p>Contributions are welcome — bug reports, feature suggestions, and pull requests alike. Please read <span class="highlight">CONTRIBUTING.md</span> before submitting. It covers how to run the test suite locally and the PR process.</p>
    </div>
  </div>

  <!-- ═══ FOOTER ═══ -->
  <div class="footer">
    <div class="footer-brand">GoQueue++ 🚀</div>
    <div class="footer-sub">Fault-Tolerant Distributed Job Processing Engine · v1.0</div>
    <div class="footer-chips">
      <span class="chip">Go</span>
      <span class="chip">Redis</span>
      <span class="chip">PostgreSQL</span>
      <span class="chip">Prometheus</span>
      <span class="chip">Grafana</span>
      <span class="chip">Docker</span>
      <span class="chip">k6</span>
    </div>
    <div style="margin-top:20px;font-size:11px;color:var(--text3);">Licensed under the terms in LICENSE · <a href="https://github.com/bhusaremayur/goqueue" style="color:var(--blue);text-decoration:none;">github.com/bhusaremayur/goqueue</a></div>
  </div>

</div>

<script>
  // ── INTERSECTION OBSERVER for section reveals ──
  const sections = document.querySelectorAll('.section');
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        e.target.classList.add('visible');
        observer.unobserve(e.target);
      }
    });
  }, { threshold: 0.08 });
  sections.forEach(s => observer.observe(s));

  // ── ANIMATED COUNTER ──
  function animateCounter(el, target, duration) {
    let start = 0;
    const step = target / (duration / 16);
    const timer = setInterval(() => {
      start += step;
      if (start >= target) {
        start = target;
        clearInterval(timer);
      }
      el.textContent = Math.floor(start).toLocaleString();
    }, 16);
  }

  const counterObserver = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        document.querySelectorAll('[data-count]').forEach(el => {
          animateCounter(el, parseInt(el.dataset.count), 1800);
        });
        counterObserver.disconnect();
      }
    });
  }, { threshold: 0.3 });

  const heroStats = document.querySelector('.hero-stats');
  if (heroStats) counterObserver.observe(heroStats);

  // ── COPY BUTTON ──
  function copyCode(btn, text) {
    navigator.clipboard.writeText(text || btn.closest('.code-block').querySelector('.code-body').innerText).then(() => {
      btn.textContent = 'copied!';
      btn.classList.add('copied');
      setTimeout(() => {
        btn.textContent = 'copy';
        btn.classList.remove('copied');
      }, 2000);
    });
  }

  // ── LOG REPLAY ──
  function replayLogs() {
    const lines = document.querySelectorAll('#logDemo .log-line');
    lines.forEach((l, i) => {
      l.style.opacity = '0';
      l.style.animation = 'none';
      setTimeout(() => {
        l.style.animation = `logAppear 0.4s ease ${i * 0.4}s forwards`;
      }, 20);
    });
  }

  document.getElementById('logDemo').addEventListener('click', replayLogs);
</script>
</body>
</html>