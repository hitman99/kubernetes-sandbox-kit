import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import LoadingButton from '@mui/lab/LoadingButton';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import kubelogo from '../static/kubelogo.png'
import {register} from "../utils/api";
import {useState} from "react";
import Alert from '@mui/material/Alert';

function Copyright(props) {
    return (
        <Typography variant="body2" color="text.secondary" align="center" {...props}>
            {'Copyright © '}
            <Link color="inherit" href="https://mui.com/">
                Your Website
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
        </Typography>
    );
}

const theme = createTheme();



export default function Register(props) {

    const [state, setState] = useState({isError: false, isLoading: false, errorMessage: ''});
    const persistState = regData => {
        localStorage.setItem("registrationData", JSON.stringify(regData))
    }

     const submit = async event => {
        event.preventDefault();
        const data = new FormData(event.currentTarget);

        setState( {isLoading: true, isError: false, errorMessage: ""});
        let res = await register({user: { email: data.get('email')}});
        try {
            if (res.success !== true) {
                setState({isError: true, isLoading: false, errorMessage: res.message});
            } else {
                setState({isError: false, isLoading: false, errorMessage: ''});
                persistState(res);
                props.setRegData(res);
            }
        } catch (err) {
            setState({isError: true, isLoading: false, errorMessage: "something went wrong, please try again"});
        }
    }
    let alert;
    if( state.isError ) {
        alert = <Alert severity="error">{state.errorMessage}</Alert>
    }
    return (
        <ThemeProvider theme={theme}>
            <Container component="main" maxWidth="xs">
                <CssBaseline />
                <Box
                    sx={{
                        marginTop: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                    }}
                >
                    <Avatar sx={{ width: 50, height: 50 }} variant="square" alt="Kubernetes" src={kubelogo} />
                    {/*<Avatar  src="static/kubelogo.png" />*/}
                    <Typography component="h1" variant="h5">
                        Kubernetes Sandbox
                    </Typography>
                    <Box component="form" onSubmit={submit} noValidate sx={{ mt: 1 }}>
                        <TextField
                            margin="normal"
                            required
                            fullWidth
                            id="email"
                            label="Email Address"
                            name="email"
                            autoComplete="email"
                            autoFocus
                        />
                        { alert }
                        <LoadingButton
                            type="submit"
                            fullWidth
                            variant="contained"
                            sx={{ mt: 3, mb: 2 }}
                            loading={state.isLoading}
                        >
                            Register
                        </LoadingButton>
                    </Box>
                </Box>
                <Copyright sx={{ mt: 8, mb: 4 }} />
            </Container>
        </ThemeProvider>
    );
}