export function get<T>(url: string): Promise<T> {
    return fetch(url)
        .then((response) => {
            if (!response.ok) {
                console.error(response);
                return []; 
                // throw new Error(response.statusText);
            }
            return response.json()
        })
}