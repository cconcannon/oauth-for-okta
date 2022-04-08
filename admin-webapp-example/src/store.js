import { writable, derived } from "svelte/store"

export const isAuthenticated = writable(false)

export const user = writable({})

export const task = writable({})

export const error = writable()

export const users = writable([])

export const tasks = writable([])

export const user_tasks = derived([tasks, user], ([$tasks, $user]) => {
    let logged_in_user_tasks = [];

    if ($user && $user.email) {
        logged_in_user_tasks = $tasks.filter((task) => task.user === $user.email);
    }

    return logged_in_user_tasks;
});