package reporting

import (
	"html/template"
	"io"
	"strings"
)

var reportTemplate = template.Must(
	template.New("report").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Assessment.Name}} - RedLab Report</title>
	<style>
		:root {
			--ink: #172033;
			--muted: #65718a;
			--line: #e5eaf2;
			--soft: #f5f7fb;
			--navy: #101a33;
			--blue: #3f5ee8;
			--red: #d9364f;
			--orange: #c76c16;
			--green: #198754;
		}

		* { box-sizing: border-box; }
		@page { size: A4; margin: 16mm; }
		body {
			margin: 0;
			color: var(--ink);
			background:
				radial-gradient(circle at 10% 0%, rgba(63, 94, 232, .12), transparent 28rem),
				#edf1f7;
			font-family: Inter, ui-sans-serif, system-ui, -apple-system, "Segoe UI", sans-serif;
			line-height: 1.55;
		}
		.page {
			width: min(1120px, calc(100% - 40px));
			margin: 34px auto;
			background: white;
			box-shadow: 0 24px 70px rgba(16, 26, 51, 0.14);
		}
		.cover {
			position: relative;
			overflow: hidden;
			min-height: 440px;
			display: flex;
			flex-direction: column;
			justify-content: center;
			padding: 64px;
			color: white;
			background: linear-gradient(130deg, #0c1732, #1f3270 70%, #3f5ee8);
		}
		.cover::before {
			position: absolute;
			inset: 0;
			background-image: linear-gradient(rgba(255,255,255,.045) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.045) 1px, transparent 1px);
			background-size: 34px 34px;
			mask-image: linear-gradient(135deg, black, transparent 72%);
			content: "";
		}
		.cover > * { position: relative; z-index: 1; }
		.cover::after {
			position: absolute;
			right: -100px;
			bottom: -150px;
			width: 390px;
			height: 390px;
			border: 1px solid rgba(255,255,255,.16);
			border-radius: 50%;
			box-shadow: 0 0 0 35px rgba(255,255,255,.035), 0 0 0 70px rgba(255,255,255,.025);
			content: "";
		}
		.brand { color: #cbd5ff; font-size: 13px; font-weight: 800; letter-spacing: .2em; text-transform: uppercase; }
		.cover h1 { max-width: 780px; margin: 22px 0 12px; font-size: clamp(30px, 5vw, 54px); line-height: 1.08; letter-spacing: -.04em; text-wrap: balance; }
		.cover p { max-width: 720px; margin: 0; color: #d9e1ff; font-size: 16px; }
		.cover-meta { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 34px; }
		.cover-pill { padding: 8px 12px; border: 1px solid rgba(255,255,255,.2); border-radius: 999px; color: #eef1ff; background: rgba(255,255,255,.09); font-size: 12px; font-weight: 750; }
		.content { padding: 52px 64px 64px; background: linear-gradient(180deg, #fff 0%, #fbfcff 100%); }
		.section { margin-top: 38px; }
		.section:first-child { margin-top: 0; }
		.section-heading { display: flex; align-items: end; justify-content: space-between; gap: 20px; margin-bottom: 16px; padding-bottom: 2px; border-bottom: 2px solid var(--ink); }
		.section-heading h2 { margin: 0 0 9px; font-size: 20px; letter-spacing: -.02em; }
		.section-heading span { margin-bottom: 9px; color: var(--muted); font-size: 12px; font-weight: 800; text-transform: uppercase; letter-spacing: .1em; }
		.overview { display: grid; grid-template-columns: 1.5fr 1fr; gap: 28px; }
		.overview p { margin: 0; color: #4e5a70; }
		.details { display: grid; grid-template-columns: 1fr 1fr; gap: 1px; overflow: hidden; border: 1px solid var(--line); border-radius: 12px; background: var(--line); }
		.detail { padding: 13px 15px; background: white; }
		.detail span { display: block; color: var(--muted); font-size: 11px; font-weight: 800; text-transform: uppercase; letter-spacing: .08em; }
		.detail strong { display: block; margin-top: 4px; font-size: 12px; word-break: break-word; }
		.summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; margin-top: 28px; }
		.summary-card { position: relative; overflow: hidden; padding: 20px; border: 1px solid var(--line); border-radius: 12px; background: var(--soft); }
		.summary-card::after { position: absolute; right: -18px; bottom: -28px; width: 80px; height: 80px; border: 12px solid rgba(63,94,232,.07); border-radius: 50%; content: ""; }
		.summary-card:nth-child(2) { background: #fff3f5; border-color: #ffd7de; }
		.summary-card:nth-child(3) { background: #fff7ed; border-color: #ffe1bf; }
		.summary-card strong { display: block; font-size: 28px; line-height: 1; }
		.summary-card span { display: block; margin-top: 7px; color: var(--muted); font-size: 12px; font-weight: 700; }
		.finding { margin: 18px 0; padding: 25px 26px; border: 1px solid var(--line); border-left: 5px solid var(--blue); border-radius: 12px; background: white; box-shadow: 0 8px 22px rgba(16,26,51,.045); page-break-inside: avoid; }
		.finding.critical { border-left-color: var(--red); }
		.finding.high { border-left-color: #e85b2a; }
		.finding.medium { border-left-color: var(--orange); }
		.finding.low { border-left-color: var(--green); }
		.finding-header { display: flex; justify-content: space-between; align-items: start; gap: 18px; }
		.finding h3 { margin: 0; font-size: 17px; }
		.severity { flex: 0 0 auto; padding: 5px 9px; border-radius: 999px; color: #3150c7; background: #edf0ff; font-size: 11px; font-weight: 900; letter-spacing: .08em; text-transform: uppercase; }
		.critical .severity { color: var(--red); background: #fff0f2; }
		.high .severity { color: #c0441c; background: #fff1eb; }
		.medium .severity { color: var(--orange); background: #fff6e8; }
		.finding p { margin: 15px 0 0; color: #4e5a70; }
		.finding-meta { display: flex; flex-wrap: wrap; gap: 8px 18px; margin-top: 17px; padding-top: 14px; border-top: 1px solid var(--line); color: var(--muted); font-size: 12px; }
		.finding-meta strong { color: var(--ink); }
		.recommendation { margin-top: 16px; padding: 13px 15px; border-radius: 8px; background: #f2f7ff; color: #3c557e; font-size: 13px; }
		.recommendation strong { color: var(--blue); }
		.finding strong { font-weight: 800; }
		.empty { padding: 25px; border: 1px dashed #bdc7d8; border-radius: 12px; color: var(--muted); text-align: center; }
		.footer { display: flex; justify-content: space-between; gap: 20px; margin-top: 48px; padding-top: 18px; border-top: 1px solid var(--line); color: var(--muted); font-size: 12px; }
		@media (max-width: 720px) {
			.page { width: 100%; margin: 0; }
			.cover, .content { padding: 34px 24px; }
			.overview, .summary { grid-template-columns: 1fr; }
			.footer { flex-direction: column; }
		}
		@media print {
			body { background: white; }
			.page { width: 100%; margin: 0; box-shadow: none; }
			.cover { min-height: 255mm; print-color-adjust: exact; -webkit-print-color-adjust: exact; }
			.finding, .summary-card { print-color-adjust: exact; -webkit-print-color-adjust: exact; }
			.section { break-inside: avoid; }
		}
	</style>
</head>
<body>
	<main class="page">
		<header class="cover">
			<div class="brand">RedLab · Security Assessment Report</div>
			<h1>{{.Assessment.Name}}</h1>
			<p>{{if .Assessment.Description}}{{.Assessment.Description}}{{else}}Authorized Active Directory security assessment{{end}}</p>
			<div class="cover-meta">
				<span class="cover-pill">Status: {{.Assessment.Status}}</span>
				<span class="cover-pill">Generated {{.GeneratedAt.Format "02 Jan 2006 · 15:04"}}</span>
			</div>
		</header>

		<div class="content">
			<section class="section">
				<div class="section-heading"><h2>Assessment overview</h2><span>Executive context</span></div>
				<div class="overview">
					<p>{{if .Assessment.Description}}{{.Assessment.Description}}{{else}}This report summarizes the evidence and security findings collected during the assessment.{{end}}</p>
					<div class="details">
						<div class="detail"><span>Assessment ID</span><strong>{{.Assessment.ID}}</strong></div>
						<div class="detail"><span>Agent</span><strong>{{.Assessment.AgentID}}</strong></div>
						<div class="detail"><span>Created</span><strong>{{.Assessment.CreatedAt.Format "02 Jan 2006 · 15:04"}}</strong></div>
						<div class="detail"><span>Updated</span><strong>{{.Assessment.UpdatedAt.Format "02 Jan 2006 · 15:04"}}</strong></div>
					</div>
				</div>
				<div class="summary">
					<div class="summary-card"><strong>{{len .Findings}}</strong><span>Total findings</span></div>
					<div class="summary-card"><strong>{{index .SeverityCounts "CRITICAL"}}</strong><span>Critical findings flagged</span></div>
					<div class="summary-card"><strong>{{index .SeverityCounts "HIGH"}}</strong><span>High findings flagged</span></div>
				</div>
			</section>

			<section class="section">
				<div class="section-heading"><h2>Security findings</h2><span>Prioritized observations</span></div>
				{{if .Findings}}
					{{range .Findings}}
					<article class="finding {{lower .Severity}}">
						<div class="finding-header">
							<h3>{{.Title}}</h3>
							<span class="severity">{{.Severity}}</span>
						</div>
						<p>{{.Description}}</p>
						<div class="finding-meta">
							<span>MITRE ATT&CK <strong>{{.TechniqueID}}</strong></span>
							<span>{{.TechniqueName}}</span>
							<span>Observed {{.CreatedAt.Format "02 Jan 2006"}}</span>
						</div>
						{{if .Recommendation}}<div class="recommendation"><strong>Recommended action:</strong> {{.Recommendation}}</div>{{end}}
					</article>
					{{end}}
				{{else}}
					<div class="empty">No findings have been recorded for this assessment.</div>
				{{end}}
			</section>

			<footer class="footer">Generated by RedLab · Active Directory Offensive Security Assessment Platform · Use only for authorized assessments.</footer>
		</div>
	</main>
</body>
</html>
`),
)

func RenderReport(writer io.Writer, report *Report) error {
	return reportTemplate.Execute(writer, report)
}
