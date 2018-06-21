export function get<T>(url: string): Promise<T> {
    return fetch(url)
        .then((response) => {
            if (!response.ok) {
                return []; 
                // throw new Error(response.statusText);
            }
            return response.json()
        })
}