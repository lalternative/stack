import { configureCoreClient, getCoreAPI } from "@app/sdk";

// Generated routes are relative to the API basePath (/api/v1). In dev Vite
// serves same-origin; override the origin via env in prod.
const origin = import.meta.env["VITE_CORE_API_URL"]?.trim().replace(/\/+$/, "") || "";
configureCoreClient({ baseURL: `${origin}/api/v1` });

export const coreAPI = getCoreAPI();
