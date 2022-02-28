const register = async (data) => {
    try {
        let res = await fetch('/register', {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'POST',
            body: JSON.stringify(data),
            cache: 'no-cache',
            //credentials: 'same-origin'
        });
        if (res.status > 400) {
            return new Error(res.statusText)
        }
        return await res.json();
    } catch(err) {
        return err;
    }
};

const instructions = async () => {
    try {
        let res = await fetch('/instructions', {
            headers: {
                'Content-Type': 'application/json'
            },
            method: 'GET',
            cache: 'no-cache'
        });
        if (res.status > 400) {
            return new Error(res.statusText)
        }
        return await res.json();
    } catch(err) {
        return err;
    }
}

export {
    register,
    instructions
}