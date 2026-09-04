const API_BASE =
    import.meta.env.VITE_API_BASE_URL ||
    "http://127.0.0.1:8080/api/v1";


async function request(path, options = {}) {
    const headers = {
        Accept: "application/json",
        ...(options.headers || {})
    };

    if (options.token) {
        headers.Authorization =
            `Bearer ${options.token}`;
    }

    if (
        options.body &&
        !headers["Content-Type"]
    ) {
        headers["Content-Type"] =
            "application/json";
    }

    let response;

    try {
        response = await fetch(
            `${API_BASE}${path}`,
            {
                ...options,
                headers
            }
        );
    } catch (error) {
        throw new Error(
            `Unable to reach RedLab API at ${API_BASE}: ${error.message}`
        );
    }

    if (response.status === 204) {
        return null;
    }

    const contentType =
        response.headers.get(
            "content-type"
        ) || "";

    if (!response.ok) {
        let message =
            `HTTP ${response.status}`;

        if (
            contentType.includes(
                "application/json"
            )
        ) {
            try {
                const body =
                    await response.json();

                message =
                    body.error ||
                    body.message ||
                    message;
            } catch {
                // Keep HTTP status message.
            }
        } else {
            try {
                const body =
                    await response.text();

                if (body) {
                    message = body;
                }
            } catch {
                // Keep HTTP status message.
            }
        }

        throw new Error(message);
    }

    if (
        contentType.includes(
            "application/json"
        )
    ) {
        return response.json();
    }

    return response.text();
}


/* =========================================================
 * AUTHENTICATION
 * ========================================================= */

export async function login(
    username,
    password
) {
    const response =
        await request(
            "/auth/login",
            {
                method: "POST",
                body: JSON.stringify({
                    username,
                    password
                })
            }
        );

    if (
        !response ||
        !response.token
    ) {
        throw new Error(
            "Login response did not contain a token."
        );
    }

    return response.token;
}


/* =========================================================
 * HEALTH
 * ========================================================= */

export async function getHealth() {
    return request(
        "/health",
        {
            method: "GET"
        }
    );
}


/* =========================================================
 * PROFILE
 * ========================================================= */

export async function getProfile(
    token
) {
    /*
     * Backend route:
     * GET /api/v1/profile
     */
    const response =
        await request(
            "/profile",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        response;
}


/* =========================================================
 * PLATFORM USER STATISTICS
 *
 * Not implemented by the current backend.
 * Kept as a compatibility export.
 * ========================================================= */

export async function getUserStats() {
    return null;
}


/* =========================================================
 * AGENTS
 * ========================================================= */

export async function getAgents(
    token
) {
    const response =
        await request(
            "/agents",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        [];
}


/* =========================================================
 * ASSESSMENTS
 * ========================================================= */

export async function getAssessments(
    token
) {
    const response =
        await request(
            "/assessments",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        [];
}


/* =========================================================
 * FINDINGS
 * ========================================================= */

export async function getFindings(
    token
) {
    const response =
        await request(
            "/findings",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        [];
}


/* =========================================================
 * TASKS
 * ========================================================= */

export async function getTasks(
    token
) {
    const response =
        await request(
            "/tasks",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        [];
}


/* =========================================================
 * CREATE TASK
 * ========================================================= */

export async function createTask(
    token,
    payload
) {
    const response =
        await request(
            "/tasks",
            {
                method: "POST",
                token,
                body: JSON.stringify(
                    payload
                )
            }
        );

    return response?.data ||
        response;
}


/* =========================================================
 * CANCEL PENDING TASK
 * ========================================================= */

export async function deleteTask(
    token,
    taskId
) {
    return request(
        `/tasks/${encodeURIComponent(
            taskId
        )}`,
        {
            method: "DELETE",
            token
        }
    );
}

export async function reviewTask(token, taskId, status) {
    return request(
        `/tasks/${encodeURIComponent(taskId)}/review`,
        {
            method: "PATCH",
            token,
            body: JSON.stringify({ status })
        }
    );
}

export async function getDirectoryUsers(token, search = "") {
    const query = search
        ? `?search=${encodeURIComponent(search)}`
        : "";
    const response = await request(
        `/directory/users${query}`,
        {
            method: "GET",
            token
        }
    );
    return response?.data || [];
}


/* =========================================================
 * DASHBOARD
 * ========================================================= */

export async function getDashboardStats(
    token
) {
    const response =
        await request(
            "/dashboard/stats",
            {
                method: "GET",
                token
            }
        );

    return response?.data ||
        response;
}


/* =========================================================
 * ASSESSMENT STATUS
 * ========================================================= */

export async function updateAssessmentStatus(
    token,
    assessmentId,
    status
) {
    const response =
        await request(
            `/assessments/${encodeURIComponent(
                assessmentId
            )}/status`,
            {
                method: "PATCH",
                token,
                body: JSON.stringify({
                    status
                })
            }
        );

    return response?.data ||
        response;
}


/* =========================================================
 * REPORT
 * ========================================================= */

export async function getReport(
    token,
    assessmentId
) {
    const headers = {
        Accept: "text/html"
    };

    if (token) {
        headers.Authorization =
            `Bearer ${token}`;
    }

    let response;

    try {
        response =
            await fetch(
                `${API_BASE}/assessments/` +
                `${encodeURIComponent(
                    assessmentId
                )}/report`,
                {
                    method: "GET",
                    headers
                }
            );
    } catch (error) {
        throw new Error(
            `Unable to reach RedLab API: ${error.message}`
        );
    }

    if (!response.ok) {
        throw new Error(
            `Failed to retrieve report: HTTP ${response.status}`
        );
    }

    return response.text();
}