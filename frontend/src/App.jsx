import React, {
    useEffect,
    useMemo,
    useState
} from "react";

import {
    createTask,
    deleteTask,
    getDirectoryUsers,
    getAgents,
    getAssessments,
    getDashboardStats,
    getFindings,
    getHealth,
    getProfile,
    getReport,
    getTasks,
    login,
    reviewTask,
    updateAssessmentStatus
} from "./services/api";
const NAV = [
    ["dashboard", "Dashboard"],
    ["agents", "Agents"],
    ["modules", "Assessment Modules"],
    ["tasks", "Tasks"],
    ["assessments", "Assessments"],
    ["findings", "Findings"],
    ["reports", "Reports"],
    ["users", "Users"]
];

const MODULES = [
    {
        id: "DOMAIN_INFO",
        module: "domain_info",
        category: "AD Discovery",
        protocol: "LDAP / RootDSE",
        technique: "T1482",
        techniqueName:
            "Domain Trust Discovery",
        description:
            "Collects domain and forest naming context information using LDAP RootDSE."
    },
    {
        id: "AD_USER_ENUMERATION",
        module: "ad_user_enumeration",
        category: "Identity Discovery",
        protocol: "LDAP",
        technique: "T1087.002",
        techniqueName:
            "Domain Account Discovery",
        description:
            "Enumerates domain user objects and account metadata."
    },
    {
        id: "AD_GROUP_ENUMERATION",
        module: "ad_group_enumeration",
        category: "Identity Discovery",
        protocol: "LDAP",
        technique: "T1069.002",
        techniqueName:
            "Permission Groups Discovery",
        description:
            "Enumerates domain groups and directory metadata."
    },
    {
        id: "AD_COMPUTER_ENUMERATION",
        module: "ad_computer_enumeration",
        category: "Host Discovery",
        protocol: "LDAP",
        technique: "T1018",
        techniqueName:
            "Remote System Discovery",
        description:
            "Enumerates domain computers, hostnames and OS information."
    },
    {
        id: "SPN_ENUMERATION",
        module: "spn_enumeration",
        category: "Kerberos Discovery",
        protocol: "LDAP",
        technique: "T1558.003",
        techniqueName:
            "Kerberoasting",
        description:
            "Finds accounts with service principal names. This is SPN discovery and Kerberoasting preparation, not ticket extraction or password cracking."
    }
];

function App() {
    const [token, setToken] =
        useState(
            localStorage.getItem(
                "redlab_token"
            ) || ""
        );

    const [page, setPage] =
        useState("dashboard");

    const [theme, setTheme] =
        useState(
            localStorage.getItem(
                "redlab_theme"
            ) || "dark"
        );

    const [username, setUsername] =
        useState("admin");

    const [password, setPassword] =
        useState("Admin123!");

    const [health, setHealth] =
        useState(null);

    const [profile, setProfile] =
        useState(null);

    const [agents, setAgents] =
        useState([]);

    const [assessments, setAssessments] =
        useState([]);

    const [findings, setFindings] =
        useState([]);

    const [userStats, setUserStats] =
        useState(null);

    const [directoryUsers, setDirectoryUsers] =
        useState([]);

    const [directoryError, setDirectoryError] =
        useState("");

   const [taskHistory, setTaskHistory] =
    useState([]);

const [dashboard, setDashboard] =
    useState(null);

    const [report, setReport] =
        useState(null);

    const [loading, setLoading] =
        useState(false);

    const [error, setError] =
        useState("");

    const [message, setMessage] =
        useState("");

    const [lastUpdated, setLastUpdated] =
        useState(null);

    const [taskAgent, setTaskAgent] =
        useState("");

    const [taskType, setTaskType] =
        useState("DOMAIN_INFO");

    const [taskPriority, setTaskPriority] =
        useState(10);

    const [taskPayload, setTaskPayload] =
        useState("{}");

    useEffect(() => {
        document.documentElement.dataset.theme =
            theme;

        localStorage.setItem(
            "redlab_theme",
            theme
        );
    }, [theme]);

    const stats = useMemo(
    () => ({
        agents:
            dashboard?.agents?.total ??
            agents.length,

        assessments:
            dashboard?.assessments?.total ??
            assessments.length,

        findings:
            dashboard?.findings?.total ??
            findings.length,

        modules:
            MODULES.length,

        users: "—",

        activeUsers: "—"
    }),
    [
        dashboard,
        agents,
        assessments,
        findings
    ]
);

    async function loadData() {
        if (!token) {
            return;
        }

        setLoading(true);
        setError("");

        try {
          const [
    healthData,
    profileData,
    agentData,
    assessmentData,
    findingData,
    taskData,
    dashboardData,
    directoryData
] = await Promise.all([
    getHealth(),
    getProfile(token),
    getAgents(token),
    getAssessments(token),
    getFindings(token),
    getTasks(token),
    getDashboardStats(token),
    getDirectoryUsers(token).catch((err) => {
        setDirectoryError(
            err.message ||
            "Active Directory is unavailable"
        );
        return [];
    })
]);

            setHealth(healthData);
            setProfile(profileData);
            setAgents(agentData);
            setAssessments(
                assessmentData
            );
            setTaskHistory(
    taskData
);

setDashboard(
    dashboardData
);
            setDirectoryUsers(
                directoryData
            );
            if (directoryData.length > 0) {
                setDirectoryError("");
            }

setLastUpdated(new Date());

if (
    agentData.length > 0 &&
    !taskAgent
            ) {
                setTaskAgent(
                    agentData[0].id
                );
            }

       } catch (err) {
            setError(
                err.message ||
                "Unable to load RedLab data"
            );
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        if (!token) {
            getHealth()
                .then(setHealth)
                .catch(() => {});
            return undefined;
        }

        loadData();
        const refreshTimer = setInterval(
            loadData,
            30000
        );

        return () =>
            clearInterval(refreshTimer);
    }, [token]);

    function logout() {
        localStorage.removeItem(
            "redlab_token"
        );

        setToken("");
        setProfile(null);
        setAgents([]);
        setAssessments([]);
        setFindings([]);
        setUserStats(null);
        setLastUpdated(null);
        setPage("dashboard");
    }

    function openReportInNewTab() {
        if (!report?.html) {
            return;
        }

        const reportBlob = new Blob(
            [report.html],
            { type: "text/html" }
        );
        const reportUrl = URL.createObjectURL(
            reportBlob
        );
        const reportWindow = window.open(
            reportUrl,
            "_blank",
            "popup,width=1200,height=900"
        );

        if (!reportWindow) {
            URL.revokeObjectURL(reportUrl);
            setError(
                "Your browser blocked the report window. Allow pop-ups to open it."
            );
            return;
        }

        window.setTimeout(
            () => URL.revokeObjectURL(reportUrl),
            60000
        );
    }

    async function handleLogin(
        event
    ) {
        event.preventDefault();

        setError("");
        setMessage("");
        setLoading(true);

        try {
            const jwt =
                await login(
                    username,
                    password
                );

            localStorage.setItem(
                "redlab_token",
                jwt
            );

            setToken(jwt);
        } catch (err) {
            setError(
                err.message ||
                "Login failed"
            );
        } finally {
            setLoading(false);
        }
    }

    function saveTaskHistory(task) {
        const next =
            [
                task,
                ...taskHistory
            ].slice(0, 25);

        setTaskHistory(next);

        localStorage.setItem(
            "redlab_task_history",
            JSON.stringify(next)
        );
    }

    async function handleCreateTask(
        event
    ) {
        event.preventDefault();

        setError("");
        setMessage("");

        if (!taskAgent) {
            setError(
                "Select an assessment agent."
            );
            return;
        }

        let payload;

        try {
            payload =
                JSON.parse(
                    taskPayload || "{}"
                );
        } catch {
            setError(
                "Payload must be valid JSON."
            );
            return;
        }

        try {
            const response =
                await createTask(
                    token,
                    {
                        agent_id:
                            taskAgent,
                        type:
                            taskType,
                        payload,
                        priority:
                            Number(
                                taskPriority
                            )
                    }
                );

            const created =
                response?.data ||
                response;

            saveTaskHistory({
                id:
                    created?.id ||
                    crypto.randomUUID(),
                agent_id:
                    taskAgent,
                type:
                    taskType,
                priority:
                    Number(
                        taskPriority
                    ),
                status:
                    created?.status ||
                    "PENDING",
                created_at:
                    created?.created_at ||
                    new Date().toISOString()
            });

            setMessage(
                `${taskType} task successfully queued.`
            );
        } catch (err) {
            setError(
                err.message ||
                "Unable to create task"
            );
        }
    }

    async function handleReport(
        assessmentId
    ) {
        setError("");

        try {
            const html =
                await getReport(
                    token,
                    assessmentId
                );

            setReport({
                assessmentId,
                html
            });
        } catch (err) {
            setError(
                err.message ||
                "Unable to load report"
            );
        }

    }

    async function handleDeleteTask(taskId) {
        setError("");
        try {
            await deleteTask(token, taskId);
            setMessage("Pending task cancelled.");
            await loadData();
        } catch (err) {
            setError(err.message || "Unable to cancel task");
        }
    }

    async function handleReviewTask(taskId, status) {
        setError("");
        try {
            await reviewTask(token, taskId, status);
            setMessage(`Task marked ${status.toLowerCase()}.`);
            await loadData();
        } catch (err) {
            setError(err.message || "Unable to review task");
        }
    }

    async function handleStatus(
        id,
        status
    ) {
        try {
            await updateAssessmentStatus(
                token,
                id,
                status
            );

            await loadData();

            setMessage(
                `Assessment status changed to ${status}.`
            );
        } catch (err) {
            setError(
                err.message ||
                "Unable to update assessment."
            );
        }
    }

    if (!token) {
        return (
            <div className="login-screen">
                <div className="login-card">
                    <Logo />

                    <div className="login-kicker">
                        OPERATOR ACCESS
                    </div>

                    <h1>
                        RedLab Console
                    </h1>

                    <p>
                        Offensive security
                        assessment control
                        center for Active
                        Directory environments.
                    </p>

                    <form
                        onSubmit={
                            handleLogin
                        }
                    >
                        <label>
                            Username
                            <input
                                value={
                                    username
                                }
                                onChange={
                                    (e) =>
                                        setUsername(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            />
                        </label>

                        <label>
                            Password
                            <input
                                type="password"
                                value={
                                    password
                                }
                                onChange={
                                    (e) =>
                                        setPassword(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            />
                        </label>

                        <button
                            className="primary-button"
                            type="submit"
                            disabled={
                                loading
                            }
                        >
                            {loading
                                ? "Authenticating..."
                                : "Sign in"}
                        </button>
                    </form>

                    {error && (
                        <div className="alert error">
                            {error}
                        </div>
                    )}

                    <div className="login-status">
                        API status:
                        <strong
                            className={
                                health?.success
                                    ? "good"
                                    : "bad"
                            }
                        >
                            {health?.success
                                ? " ONLINE"
                                : " OFFLINE"}
                        </strong>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="app-shell">
            <aside className="sidebar">
                <div className="sidebar-header">
                    <Logo compact />

                    <div className="sidebar-caption">
                        OFFENSIVE SECURITY
                    </div>
                </div>

                <nav>
                    {NAV.map(
                        ([id, label]) => (
                            <button
                                key={id}
                                className={
                                    page === id
                                        ? "nav-item active"
                                        : "nav-item"
                                }
                                onClick={() =>
                                    setPage(id)
                                }
                            >
                                <span className="nav-dot" />
                                {label}
                            </button>
                        )
                    )}
                </nav>

                <div className="sidebar-bottom">
                    <div className="operator-card">
                        <span className="operator-dot" />

                        <div>
                            <strong>
                                {
                                    profile?.username ||
                                    "Operator"
                                }
                            </strong>

                            <small>
                                {
                                    profile?.role ||
                                    "authenticated"
                                }
                            </small>
                        </div>
                    </div>

                    <button
                        className="logout-button"
                        onClick={
                            logout
                        }
                    >
                        Logout
                    </button>
                </div>
            </aside>

            <main className="main">
                <header className="topbar">
                    <div>
                        <div className="topbar-kicker">
                            REDLAB CONSOLE
                        </div>

                        <h1>
                            {
                                NAV.find(
                                    ([id]) =>
                                        id ===
                                        page
                                )?.[1]
                            }
                        </h1>

                        <p>
                            Active Directory
                            security assessment
                            workspace
                        </p>
                    </div>

                    <div className="top-actions">
                        <span
                            className={
                                health?.success
                                    ? "backend-status online"
                                    : "backend-status offline"
                            }
                        >
                            ●{" "}
                            {health?.success
                                ? "Backend Online"
                                : "Backend Offline"}
                        </span>

                        <span className="sync-status">
                            {lastUpdated
                                ? `Live sync · ${lastUpdated.toLocaleTimeString()}`
                                : "Waiting for sync"}
                        </span>

                        <button
                            className="icon-button"
                            title="Toggle light / dark mode"
                            onClick={() =>
                                setTheme(
                                    theme ===
                                        "dark"
                                        ? "light"
                                        : "dark"
                                )
                            }
                        >
                            {theme ===
                            "dark"
                                ? "☼"
                                : "☾"}
                        </button>

                        <button
                            className="refresh-button"
                            onClick={
                                loadData
                            }
                            disabled={
                                loading
                            }
                        >
                            {loading
                                ? "Refreshing"
                                : "Refresh"}
                        </button>
                    </div>
                </header>

                <section className="content">
                    {error && (
                        <div className="alert error">
                            {error}
                        </div>
                    )}

                    {message && (
                        <div className="alert success">
                            {message}
                        </div>
                    )}

                    {page ===
                        "dashboard" && (
                        <Dashboard
                            stats={stats}
                            agents={agents}
                            findings={
                                findings
                            }
                            assessments={
                                assessments
                            }
                            taskHistory={
                                taskHistory
                            }
                            onReport={
                                handleReport
                            }
                        />
                    )}

                    {page ===
                        "agents" && (
                        <AgentsPage
                            agents={
                                agents
                            }
                        />
                    )}

                    {page ===
                        "modules" && (
                        <ModulesPage />
                    )}

                    {page ===
                        "tasks" && (
                        <TasksPage
                            agents={
                                agents
                            }
                            history={
                                taskHistory
                            }
                            taskAgent={
                                taskAgent
                            }
                            setTaskAgent={
                                setTaskAgent
                            }
                            taskType={
                                taskType
                            }
                            setTaskType={
                                setTaskType
                            }
                            priority={
                                taskPriority
                            }
                            setPriority={
                                setTaskPriority
                            }
                            payload={
                                taskPayload
                            }
                            setPayload={
                                setTaskPayload
                            }
                            onSubmit={
                                handleCreateTask
                            }
                            onDelete={
                                handleDeleteTask
                            }
                            onReview={
                                handleReviewTask
                            }
                        />
                    )}

                    {page ===
                        "assessments" && (
                        <AssessmentsPage
                            assessments={
                                assessments
                            }
                            onReport={
                                handleReport
                            }
                            onStatus={
                                handleStatus
                            }
                        />
                    )}

                    {page ===
                        "findings" && (
                        <FindingsPage
                            findings={
                                findings
                            }
                        />
                    )}

                    {page ===
                        "reports" && (
                        <ReportsPage
                            assessments={
                                assessments
                            }
                            onReport={
                                handleReport
                            }
                        />
                    )}

                    {page ===
                        "users" && (
                        <UsersPage
                            stats={userStats}
                            directoryUsers={directoryUsers}
                            directoryError={directoryError}
                        />
                    )}
                </section>
            </main>

            {report && (
                <div
                    className="modal-backdrop"
                    onClick={() =>
                        setReport(null)
                    }
                >
                    <div
                        className="report-modal"
                        onClick={(e) =>
                            e.stopPropagation()
                        }
                    >
                        <div className="report-modal-header">
                            <div>
                                <span className="eyebrow">
                                    REPORT PREVIEW
                                </span>
                                <h2>
                                    Assessment
                                    Report
                                </h2>
                                <p className="report-modal-subtitle">
                                    Review the executive summary and prioritized findings before sharing.
                                </p>
                            </div>

                            <div className="report-modal-actions">
                                <span className="report-preview-status">
                                    <span className="status-dot" />
                                    Live preview
                                </span>
                                <button
                                    className="secondary-button"
                                    onClick={
                                        openReportInNewTab
                                    }
                                >
                                    Open full report
                                </button>

                            <button
                                className="icon-button"
                                onClick={() =>
                                    setReport(
                                        null
                                    )
                                }
                            >
                                ×
                            </button>
                            </div>
                        </div>

                        <div className="report-preview-canvas">
                            <div className="report-preview-toolbar">
                                <span>REDLAB / REPORT VIEWER</span>
                                <span>Generated assessment document</span>
                            </div>
                            <iframe
                                className="report-frame"
                                title="RedLab assessment report"
                                srcDoc={report.html}
                            />
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function Logo({ compact = false }) {
    return (
        <div
            className={
                compact
                    ? "logo compact"
                    : "logo"
            }
        >
            <span>Red</span>
            <strong>Lab</strong>
        </div>
    );
}

function Dashboard({
    stats,
    agents,
    findings,
    assessments,
    taskHistory,
    onReport
}) {
    return (
        <>
            <section className="hero">
                <div>
                    <div className="eyebrow">
                        OPERATOR CONSOLE
                    </div>

                    <h2>
                        Active Directory
                        Security Assessment
                    </h2>

                    <p>
                        Centralized visibility
                        across agents,
                        assessment modules,
                        findings and reports.
                    </p>
                </div>
            </section>

            <div className="metrics">
                <Metric
                    label="Agents"
                    value={
                        stats.agents
                    }
                    hint="registered endpoints"
                />

                <Metric
                    label="Assessment Modules"
                    value={
                        stats.modules
                    }
                    hint="implemented AD modules"
                />

                <Metric
                    label="Assessments"
                    value={
                        stats.assessments
                    }
                    hint="assessment operations"
                />

                <Metric
                    label="Findings"
                    value={
                        stats.findings
                    }
                    hint="security observations"
                />

                <Metric
                    label="Platform Users"
                    value={
                        stats.users
                    }
                    hint={
                        stats.users ===
                        "—"
                            ? "backend API pending"
                            : "registered users"
                    }
                />

                <Metric
                    label="Active Users"
                    value={
                        stats.activeUsers
                    }
                    hint="enabled accounts"
                />
            </div>

            <div className="grid-two">
                <Panel
                    title="Assessment Agents"
                    subtitle="Registered endpoints"
                    count={
                        agents.length
                    }
                >
                    <AgentTable
                        agents={
                            agents.slice(
                                0,
                                6
                            )
                        }
                    />
                </Panel>

                <Panel
                    title="Findings"
                    subtitle="Latest security observations"
                    count={
                        findings.length
                    }
                >
                    <FindingTable
                        findings={
                            findings.slice(
                                0,
                                6
                            )
                        }
                    />
                </Panel>
            </div>

            <Panel
                title="Assessments"
                subtitle="Assessment lifecycle"
                count={
                    assessments.length
                }
            >
                <AssessmentTable
                    assessments={
                        assessments.slice(
                            0,
                            8
                        )
                    }
                    onReport={
                        onReport
                    }
                />
            </Panel>

            <Panel
                title="Live task activity"
                subtitle="Latest tasks returned by the backend"
                count={taskHistory.length}
            >
                <ActivityFeed
                    tasks={taskHistory.slice(0, 6)}
                />
            </Panel>
        </>
    );
}

function Metric({
    label,
    value,
    hint
}) {
    return (
        <div className="metric">
            <span>{label}</span>
            <strong>{value}</strong>
            <small>{hint}</small>
        </div>
    );
}

function ActivityFeed({ tasks }) {
    if (!tasks.length) {
        return (
            <Empty text="No task activity has been reported yet." />
        );
    }

    return (
        <div className="activity-feed">
            {tasks.map((task) => (
                <div
                    className="activity-item"
                    key={task.id}
                >
                    <span className="activity-pulse" />
                    <div className="data-main">
                        <strong>
                            {task.type ||
                                "Assessment task"}
                        </strong>
                        <small>
                            {task.status ||
                                "PENDING"}{" "}
                            · Agent{" "}
                            {task.agent_id ||
                                "unassigned"}
                        </small>
                    </div>
                    <span className="badge">
                        {task.priority ??
                            "—"}
                    </span>
                </div>
            ))}
        </div>
    );
}

function AgentsPage({ agents }) {
    return (
        <Panel
            title="Registered Agents"
            subtitle="Endpoints authorized to execute assessment tasks"
            count={
                agents.length
            }
        >
            {agents.length ? (
                <AgentTable
                    agents={
                        agents
                    }
                    detailed
                />
            ) : (
                <Empty text="No agents registered." />
            )}
        </Panel>
    );
}

function AgentTable({
    agents,
    detailed = false
}) {
    if (!agents.length) {
        return (
            <Empty text="No agent data available." />
        );
    }

    return (
        <div className="data-list">
            {agents.map(
                (agent) => (
                    <div
                        className="data-row"
                        key={
                            agent.id
                        }
                    >
                        <div className="data-main">
                            <strong>
                                {
                                    agent.name
                                }
                            </strong>

                            <small>
                                {
                                    agent.hostname ||
                                    "Unknown host"
                                }
                                {" · "}
                                {
                                    agent.ip_address ||
                                    "No IP"
                                }
                            </small>

                            {detailed && (
                                <small>
                                    {
                                        agent.operating_system ||
                                        "Unknown OS"
                                    }
                                    {" · "}
                                    {
                                        agent.version ||
                                        "Unknown version"
                                    }
                                </small>
                            )}
                        </div>

                        <span
                            className={
                                agent.status ===
                                "ONLINE"
                                    ? "badge good"
                                    : "badge"
                            }
                        >
                            {
                                agent.status ||
                                "UNKNOWN"
                            }
                        </span>
                    </div>
                )
            )}
        </div>
    );
}

function ModulesPage() {
    return (
        <>
            <section className="hero compact-hero">
                <div>
                    <div className="eyebrow">
                        ASSESSMENT ENGINE
                    </div>

                    <h2>
                        Active Directory
                        Modules
                    </h2>

                    <p>
                        Each module represents
                        a concrete assessment
                        capability implemented
                        by the RedLab agent.
                    </p>
                </div>
            </section>

            <div className="module-grid">
                {MODULES.map(
                    (module) => (
                        <article
                            className="module-card"
                            key={
                                module.id
                            }
                        >
                            <div className="module-top">
                                <span>
                                    {
                                        module.category
                                    }
                                </span>

                                <strong>
                                    {
                                        module.technique
                                    }
                                </strong>
                            </div>

                            <h3>
                                {
                                    module.id
                                }
                            </h3>

                            <div className="module-code">
                                {
                                    module.module
                                }
                            </div>

                            <p>
                                {
                                    module.description
                                }
                            </p>

                            <div className="module-meta">
                                <div>
                                    <span>
                                        Protocol
                                    </span>
                                    <strong>
                                        {
                                            module.protocol
                                        }
                                    </strong>
                                </div>

                                <div>
                                    <span>
                                        MITRE ATT&CK
                                    </span>
                                    <strong>
                                        {
                                            module.techniqueName
                                        }
                                    </strong>
                                </div>
                            </div>
                        </article>
                    )
                )}
            </div>
        </>
    );
}

function TasksPage({
    agents,
    history,
    taskAgent,
    setTaskAgent,
    taskType,
    setTaskType,
    priority,
    setPriority,
    payload,
    setPayload,
    onSubmit,
    onDelete,
    onReview
}) {
    return (
        <div className="grid-two">
            <Panel
                title="Dispatch Assessment Task"
                subtitle="Create a module execution for an agent"
            >
                {!agents.length ? (
                    <Empty text="You need a registered agent before creating a task." />
                ) : (
                    <form
                        className="task-form"
                        onSubmit={
                            onSubmit
                        }
                    >
                        <label>
                            Target Agent
                            <select
                                value={
                                    taskAgent
                                }
                                onChange={
                                    (e) =>
                                        setTaskAgent(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            >
                                {agents.map(
                                    (
                                        agent
                                    ) => (
                                        <option
                                            key={
                                                agent.id
                                            }
                                            value={
                                                agent.id
                                            }
                                        >
                                            {
                                                agent.name
                                            }
                                            {" — "}
                                            {
                                                agent.hostname
                                            }
                                        </option>
                                    )
                                )}
                            </select>
                        </label>

                        <label>
                            Assessment Module
                            <select
                                value={
                                    taskType
                                }
                                onChange={
                                    (e) =>
                                        setTaskType(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            >
                                {MODULES.map(
                                    (
                                        module
                                    ) => (
                                        <option
                                            key={
                                                module.id
                                            }
                                            value={
                                                module.id
                                            }
                                        >
                                            {
                                                module.id
                                            }
                                        </option>
                                    )
                                )}
                            </select>
                        </label>

                        <label>
                            Priority
                            <input
                                type="number"
                                min="0"
                                max="100"
                                value={
                                    priority
                                }
                                onChange={
                                    (e) =>
                                        setPriority(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            />
                        </label>

                        <label>
                            Module Payload
                            <textarea
                                rows="8"
                                value={
                                    payload
                                }
                                onChange={
                                    (e) =>
                                        setPayload(
                                            e
                                                .target
                                                .value
                                        )
                                }
                            />
                        </label>

                        <button
                            className="primary-button"
                            type="submit"
                        >
                            Dispatch Task
                        </button>
                    </form>
                )}
            </Panel>

            <Panel
                title="Task Activity"
                subtitle="Tasks created during this operator session"
                count={
                    history.length
                }
            >
                {history.length ? (
                    <div className="data-list">
                        {history.map(
                            (task) => (
                                <div
                                    className="data-row"
                                    key={
                                        task.id
                                    }
                                >
                                    <div className="data-main">
                                        <strong>
                                            {
                                                task.type
                                            }
                                        </strong>

                                        <small>
                                            Agent:{" "}
                                            {
                                                task.agent_id
                                            }
                                        </small>

                                        <small>
                                            Priority:{" "}
                                            {
                                                task.priority
                                            }
                                        </small>
                                    </div>

                                    <span className="badge">
                                        {
                                            task.status
                                        }
                                    </span>
                                    {task.status ===
                                        "PENDING" && (
                                        <button
                                            className="small-button danger-button"
                                            onClick={() =>
                                                onDelete(
                                                    task.id
                                                )
                                            }
                                        >
                                            Cancel
                                        </button>
                                    )}
                                    {(task.status ===
                                        "COMPLETED" ||
                                        task.status ===
                                            "FAILED") &&
                                        task.review_status !==
                                            "APPROVED" &&
                                        task.review_status !==
                                            "REJECTED" && (
                                        <>
                                            <button
                                                className="small-button"
                                                onClick={() =>
                                                    onReview(
                                                        task.id,
                                                        "APPROVED"
                                                    )
                                                }
                                            >
                                                Approve
                                            </button>
                                            <button
                                                className="small-button danger-button"
                                                onClick={() =>
                                                    onReview(
                                                        task.id,
                                                        "REJECTED"
                                                    )
                                                }
                                            >
                                                Reject
                                            </button>
                                        </>
                                    )}
                                    {task.review_status &&
                                        task.review_status !==
                                            "PENDING" && (
                                        <span className="badge good">
                                            {task.review_status}
                                        </span>
                                    )}
                                </div>
                            )
                        )}
                    </div>
                ) : (
                    <Empty text="No tasks created in this browser session." />
                )}
            </Panel>
        </div>
    );
}

function AssessmentsPage({
    assessments,
    onReport,
    onStatus
}) {
    return (
        <Panel
            title="Assessments"
            subtitle="Assessment lifecycle and report access"
            count={
                assessments.length
            }
        >
            <AssessmentTable
                assessments={
                    assessments
                }
                detailed
                onReport={
                    onReport
                }
                onStatus={
                    onStatus
                }
            />
        </Panel>
    );
}

function AssessmentTable({
    assessments,
    onReport,
    onStatus,
    detailed = false
}) {
    if (!assessments.length) {
        return (
            <Empty text="No assessments available." />
        );
    }

    return (
        <div className="data-list">
            {assessments.map(
                (assessment) => (
                    <div
                        className="data-row assessment-row"
                        key={
                            assessment.id
                        }
                    >
                        <div className="data-main">
                            <strong>
                                {
                                    assessment.name
                                }
                            </strong>

                            <small>
                                {
                                    assessment.description ||
                                    "No description"
                                }
                            </small>

                            {detailed && (
                                <small>
                                    Agent:{" "}
                                    {
                                        assessment.agent_id
                                    }
                                </small>
                            )}
                        </div>

                        <div className="row-actions">
                            <span className="badge">
                                {
                                    assessment.status
                                }
                            </span>

                            <button
                                className="small-button"
                                onClick={() =>
                                    onReport(
                                        assessment.id
                                    )
                                }
                            >
                                View Report
                            </button>

                            {onStatus &&
                                assessment.status !==
                                    "COMPLETED" && (
                                    <button
                                        className="small-button"
                                        onClick={() =>
                                            onStatus(
                                                assessment.id,
                                                "COMPLETED"
                                            )
                                        }
                                    >
                                        Complete
                                    </button>
                                )}
                        </div>
                    </div>
                )
            )}
        </div>
    );
}

function FindingsPage({
    findings
}) {
    return (
        <Panel
            title="Security Findings"
            subtitle="Assessment observations mapped to ATT&CK techniques"
            count={
                findings.length
            }
        >
            {findings.length ? (
                <FindingTable
                    findings={
                        findings
                    }
                    detailed
                />
            ) : (
                <Empty text="No findings available." />
            )}
        </Panel>
    );
}

function FindingTable({
    findings,
    detailed = false
}) {
    if (!findings.length) {
        return (
            <Empty text="No findings available." />
        );
    }

    return (
        <div className="data-list">
            {findings.map(
                (finding) => (
                    <div
                        className="data-row"
                        key={
                            finding.id
                        }
                    >
                        <div className="data-main">
                            <strong>
                                {
                                    finding.title
                                }
                            </strong>

                            <small>
                                {
                                    finding.technique_id
                                }
                                {" · "}
                                {
                                    finding.technique_name
                                }
                            </small>

                            {detailed &&
                                finding.recommendation && (
                                    <small>
                                        {
                                            finding.recommendation
                                        }
                                    </small>
                                )}
                        </div>

                        <span className="badge">
                            {
                                finding.severity ||
                                "INFO"
                            }
                        </span>
                    </div>
                )
            )}
        </div>
    );
}

function ReportsPage({
    assessments,
    onReport
}) {
    return (
        <Panel
            title="Assessment Reports"
            subtitle="Generate and review assessment output"
            count={
                assessments.length
            }
        >
            {!assessments.length ? (
                <Empty text="No completed or active assessments available." />
            ) : (
                <div className="report-grid">
                    {assessments.map(
                        (assessment) => (
                            <article
                                className="report-card"
                                key={
                                    assessment.id
                                }
                            >
                                <div>
                                    <div className="eyebrow">
                                        ASSESSMENT
                                    </div>

                                    <h3>
                                        {
                                            assessment.name
                                        }
                                    </h3>

                                    <p>
                                        {
                                            assessment.description ||
                                            "Security assessment report"
                                        }
                                    </p>

                                    <span className="badge">
                                        {
                                            assessment.status
                                        }
                                    </span>
                                </div>

                                <button
                                    className="primary-button"
                                    onClick={() =>
                                        onReport(
                                            assessment.id
                                        )
                                    }
                                >
                                    Preview Report
                                </button>
                            </article>
                        )
                    )}
                </div>
            )}
        </Panel>
    );
}

function UsersPage({
    stats,
    directoryUsers,
    directoryError
}) {
    const enabledUsers = directoryUsers.filter(
        (user) => user.enabled
    ).length;

    return (
        <>
            <div className="metrics user-metrics">
                <Metric
                    label="Total Users"
                    value={
                        directoryUsers.length
                    }
                    hint="Active Directory users"
                />

                <Metric
                    label="Active Users"
                    value={
                        enabledUsers
                    }
                    hint="enabled directory accounts"
                />

                <Metric
                    label="Inactive Users"
                    value={
                        directoryUsers.length -
                        enabledUsers
                    }
                    hint="disabled directory accounts"
                />

                <Metric
                    label="LDAP Users"
                    value={
                        directoryUsers.length
                    }
                    hint="synced from configured LDAP"
                />
            </div>

            <Panel
                title="Active Directory users"
                subtitle="Read-only directory view from the configured LDAP service account"
                count={directoryUsers.length}
            >
                {directoryError && (
                    <div className="alert error">
                        Active Directory could not be reached: {directoryError}
                    </div>
                )}
                {directoryUsers.length ? (
                    <div className="data-list">
                        {directoryUsers.map((user) => (
                            <div
                                className="data-row"
                                key={user.dn}
                            >
                                <div className="data-main">
                                    <strong>
                                        {user.display_name ||
                                            user.username}
                                    </strong>
                                    <small>
                                        {user.username}
                                        {user.email
                                            ? ` · ${user.email}`
                                            : ""}
                                    </small>
                                    <small>
                                        {user.groups?.length
                                            ? user.groups.join(
                                                  " · "
                                              )
                                            : "No direct group memberships"}
                                    </small>
                                </div>
                                <span
                                    className={
                                        user.enabled
                                            ? "badge good"
                                            : "badge"
                                    }
                                >
                                    {user.enabled
                                        ? "ENABLED"
                                        : "DISABLED"}
                                </span>
                            </div>
                        ))}
                    </div>
                ) : (
                    <Empty text="No directory users returned. Check LDAP connectivity and service-account permissions." />
                )}
            </Panel>
        </>
    );
}

function Panel({
    title,
    subtitle,
    count,
    children
}) {
    return (
        <section className="panel">
            <div className="panel-header">
                <div>
                    <h2>{title}</h2>

                    {subtitle && (
                        <p>
                            {subtitle}
                        </p>
                    )}
                </div>

                {count !== undefined && (
                    <span className="panel-count">
                        {count}
                    </span>
                )}
            </div>

            {children}
        </section>
    );
}

function Empty({ text }) {
    return (
        <div className="empty">
            {text}
        </div>
    );
}

export default App;