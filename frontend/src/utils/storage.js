const loadState = () => {
    let reg = localStorage.getItem("registrationData");
    if (reg) {
        return JSON.parse(reg)
    } else {
        return {
                user: {
                    email: '',
                    id: ''
                },
                kubernetes: {
                    namespace: '',
                    serverVersion: ''
                }
            }
    }
}

export {
    loadState
}