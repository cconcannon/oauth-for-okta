<script>
    import { users } from "../store"
    export let user = {}

    let isChecked;

    function userSelected() {
        console.log(isChecked)

        let updatedUsers = $users.map((user) => {
            if (currentUser.id === user.id) {
                currentUser.selected = isChecked
            }
            return currentUser
        })

        users.set(updatedUsers)
        console.log($users)
    }
</script>

<style>
    .selected {
        color: red;
        text-decoration: bold;
    }
</style>

<main>
    <li class="list-group-item">
        <input
            type="checkbox"
            class="form-check-input"
            id="exampleCheck1"
            bind:checked={isChecked}
            on:change={(e) => userSelected(e)} />
        <span class:selected={user.selected}>{user.id}</span>
    </li>
</main>