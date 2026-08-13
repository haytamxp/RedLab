const API_BASE =
    import.meta.env.VITE_API_BASE_URL ||
    "http://127.0.0.1:8080/api/v1";

async function request(path, options = {}) {
    const headers = {
        Accept: "application/json",
        ...(options.headers || {})
    };

    const token = options.token;

    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    if (options.body && !headers["Content-Type"]) {
        headers["Content-Type"] = "application/json";
    }

    const response = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers
    });

    if (response.status === 204) {
        return null;
    }

    const contentType =
        response.headers.get("content-type") || "";

    if (!response.ok) {
        let message = `HTTP ${response.status}`;

        if (contentType.includes("application/json")) {
            const body = await response.json();

            message =
                body.error ||
                body.message ||
                message;
        } else {
            const body = await response.text();

            if (body) {
                message = body;
            }
        }

        throw new Error(message);
    }

    if (contentType.includes("application/json")) {
        return response.json();
    }

    return response.text();
}

export async function login(
    username,
    password,
) {
    const response = await request(
        "/auth/login",
        {
            method: "POST",
            body: JSON.stringify({
                username,
                password
            })
        },
    );

    if (
        !response ||
        !response.token
    ) {
        throw new Error(
            "Login response did not contain a token",
        );
    }

    return response.token;
}

export async function getAgents(token) {
    const response = await request(
        "/agents",
        {
            method: "GET",
            token
        },
    );

    return response?.data || [];
}

export async function getAssessments(token) {
    const response = await request(
        "/assessments",
        {
            method: "GET",
            token
        },
    );

    return response?.data || [];
}

export async function getFindings(token) {
    const response = await request(
        "/findings",
        {
            method: "GET",
            token
        },
    );

    return response?.data || [];
}

export async function getReport(
    token,
    assessmentId,
) {
    const headers = {
        Accept: "text/html"
    };

    if (token) {
        headers.Authorization =
            `Bearer ${token}`;
    }

    const response = await fetch(
        `${API_BASE}/assessments/` +
        `${encodeURIComponent(assessmentId)}/report`,
        {
            method: "GET",
            headers
        },
    );

    if (!response.ok) {
        throw new Error(
            `Failed to retrieve report: HTTP ${response.status}`,
        );
    }

    return response.text();
}