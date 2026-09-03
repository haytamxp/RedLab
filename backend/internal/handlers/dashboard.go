package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/haytamxp/redlab/backend/internal/database"
)

type DashboardStats struct {
	Agents      AgentStats      `json:"agents"`
	Tasks       TaskStats       `json:"tasks"`
	Assessments AssessmentStats `json:"assessments"`
	Findings    FindingStats    `json:"findings"`
	MITRE       []MITREStat     `json:"mitre"`
	TaskTrend   []TaskTrendStat `json:"task_trend"`
}

type AgentStats struct {
	Total   int64 `json:"total"`
	Online  int64 `json:"online"`
	Offline int64 `json:"offline"`
}

type TaskStats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Claimed   int64 `json:"claimed"`
	Completed int64 `json:"completed"`
	Failed    int64 `json:"failed"`
}

type AssessmentStats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Running   int64 `json:"running"`
	Completed int64 `json:"completed"`
	Failed    int64 `json:"failed"`
}

type FindingStats struct {
	Total    int64 `json:"total"`
	Critical int64 `json:"critical"`
	High     int64 `json:"high"`
	Medium   int64 `json:"medium"`
	Low      int64 `json:"low"`
	Info     int64 `json:"info"`
}

type MITREStat struct {
	TechniqueID string `json:"technique_id"`
	Count       int64  `json:"count"`
}

type TaskTrendStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}


func DashboardStatsHandler(
	c *gin.Context,
) {
	ctx := c.Request.Context()

	var stats DashboardStats

	/*
	 * AGENTS
	 */

	err := database.DB.QueryRow(
		ctx,
		`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'ONLINE'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) <> 'ONLINE'
			)
		FROM agents
		`,
	).Scan(
		&stats.Agents.Total,
		&stats.Agents.Online,
		&stats.Agents.Offline,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate agent statistics",
			},
		)
		return
	}

	/*
	 * TASKS
	 */

	err = database.DB.QueryRow(
		ctx,
		`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'PENDING'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'CLAIMED'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'COMPLETED'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'FAILED'
			)
		FROM tasks
		`,
	).Scan(
		&stats.Tasks.Total,
		&stats.Tasks.Pending,
		&stats.Tasks.Claimed,
		&stats.Tasks.Completed,
		&stats.Tasks.Failed,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate task statistics",
			},
		)
		return
	}

	/*
	 * ASSESSMENTS
	 */

	err = database.DB.QueryRow(
		ctx,
		`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'PENDING'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) IN (
					'RUNNING',
					'IN_PROGRESS'
				)
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'COMPLETED'
			),
			COUNT(*) FILTER (
				WHERE UPPER(status) = 'FAILED'
			)
		FROM assessments
		`,
	).Scan(
		&stats.Assessments.Total,
		&stats.Assessments.Pending,
		&stats.Assessments.Running,
		&stats.Assessments.Completed,
		&stats.Assessments.Failed,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate assessment statistics",
			},
		)
		return
	}

	/*
	 * FINDINGS
	 */

	err = database.DB.QueryRow(
		ctx,
		`
		SELECT
			COUNT(*),
			COUNT(*) FILTER (
				WHERE UPPER(severity) = 'CRITICAL'
			),
			COUNT(*) FILTER (
				WHERE UPPER(severity) = 'HIGH'
			),
			COUNT(*) FILTER (
				WHERE UPPER(severity) = 'MEDIUM'
			),
			COUNT(*) FILTER (
				WHERE UPPER(severity) = 'LOW'
			),
			COUNT(*) FILTER (
				WHERE UPPER(severity) IN (
					'INFO',
					'INFORMATIONAL'
				)
			)
		FROM findings
		`,
	).Scan(
		&stats.Findings.Total,
		&stats.Findings.Critical,
		&stats.Findings.High,
		&stats.Findings.Medium,
		&stats.Findings.Low,
		&stats.Findings.Info,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate finding statistics",
			},
		)
		return
	}

	/*
	 * MITRE ATT&CK
	 */

	rows, err := database.DB.Query(
		ctx,
		`
		SELECT
			technique_id,
			COUNT(*) AS count
		FROM findings
		WHERE technique_id IS NOT NULL
		  AND technique_id <> ''
		GROUP BY technique_id
		ORDER BY count DESC, technique_id ASC
		LIMIT 10
		`,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate MITRE statistics",
			},
		)
		return
	}

	for rows.Next() {
		var item MITREStat

		if err := rows.Scan(
			&item.TechniqueID,
			&item.Count,
		); err != nil {
			rows.Close()

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":
						"failed to read MITRE statistics",
				},
			)
			return
		}

		stats.MITRE =
			append(
				stats.MITRE,
				item,
			)
	}

	rows.Close()

	if stats.MITRE == nil {
		stats.MITRE =
			[]MITREStat{}
	}

	/*
	 * TASK TREND — LAST 7 DAYS
	 */

	trendRows, err :=
		database.DB.Query(
			ctx,
			`
			SELECT
				TO_CHAR(
					day,
					'YYYY-MM-DD'
				),
				COUNT(t.id)
			FROM generate_series(
				CURRENT_DATE -
					INTERVAL '6 days',
				CURRENT_DATE,
				INTERVAL '1 day'
			) AS day
			LEFT JOIN tasks t
				ON t.created_at::date =
					day::date
			GROUP BY day
			ORDER BY day
			`,
		)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"success": false,
				"error":
					"failed to calculate task trend",
			},
		)
		return
	}

	for trendRows.Next() {
		var item TaskTrendStat

		if err := trendRows.Scan(
			&item.Date,
			&item.Count,
		); err != nil {
			trendRows.Close()

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"success": false,
					"error":
						"failed to read task trend",
				},
			)
			return
		}

		stats.TaskTrend =
			append(
				stats.TaskTrend,
				item,
			)
	}

	trendRows.Close()

	if stats.TaskTrend == nil {
		stats.TaskTrend =
			[]TaskTrendStat{}
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    stats,
		},
	)
}