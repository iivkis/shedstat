export function formatNumber(n: Number | String | undefined): string {
    if (n instanceof Number) {
        return n.toLocaleString('en-US')
    }
    return ""
}