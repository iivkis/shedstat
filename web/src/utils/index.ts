export function formatNumber(n: number | string | Array<number | string> | undefined): string {
    if (typeof n === 'number') {
        return n.toLocaleString('en-US')
    }
    return ""
}