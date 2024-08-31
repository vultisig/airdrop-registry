export default {
  root: "/",
  basePath: import.meta.env.VITE_BASE_PATH ? import.meta.env.VITE_BASE_PATH : "/",
  asset: "/chains/:chainKey",
  chains: "/chains",
  import: "/import",
};
