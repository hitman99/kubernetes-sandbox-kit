import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import Container from '@mui/material/Container';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Tooltip from '@mui/material/Tooltip';
import MenuItem from '@mui/material/MenuItem';
import gopher from '../static/gopher.png'
import {loadState} from "../utils/storage";
import MarkEmailReadIcon from '@mui/icons-material/MarkEmailRead';
import KeyIcon from '@mui/icons-material/Key';
import Grid from '@mui/material/Grid';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Paper from '@mui/material/Paper';
import Stack from '@mui/material/Stack';
import Instructions from "../instructions";

const settings = ['Logout'];



const ResponsiveAppBar = (props) => {
    const [anchorElNav, setAnchorElNav] = React.useState(null);
    const [anchorElUser, setAnchorElUser] = React.useState(null);

    const handleOpenNavMenu = (event) => {
        setAnchorElNav(event.currentTarget);
    };
    const handleOpenUserMenu = (event) => {
        setAnchorElUser(event.currentTarget);
    };

    const handleCloseNavMenu = () => {
        setAnchorElNav(null);
    };

    const handleCloseUserMenu = (action) => {
        if (action === 'Logout') {
            unregister()
        }
        setAnchorElUser(null);
    };

    const unregister = () => {
        localStorage.removeItem("registrationData");
        window.location = '/';
        props.setRegData(loadState())
    }

    const {user, kubernetes} = props.regData;
    return (
        <Box sx={{ flexGrow: 1 }}>

        <AppBar position="static">
            <Container maxWidth="xxl">

                <Toolbar disableGutters>
                    <Box sx={{
                        display: 'flex',
                        alignItems: 'center',
                        width: 'fit-content',
                        borderRadius: 1,
                        '& svg': {
                            m: 1.5,
                        },
                        '& hr': {
                            mx: 0.5,
                        },
                    }}>

                        <MarkEmailReadIcon /> {user.email}
                        <KeyIcon /> {user.id}
                    </Box>



                    <Typography
                        variant="h6"
                        noWrap
                        component="div"
                        sx={{ flexGrow: 1, display: { xs: 'flex', md: 'none' } }}
                    >

                    </Typography>
                    <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }}>

                    </Box>

                    <Box sx={{ flexGrow: 0 }} >
                        <Tooltip title="Open settings">
                            <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                                <Avatar alt="Remy Sharp" src={gopher} />
                            </IconButton>
                        </Tooltip>
                        <Menu
                            sx={{ mt: '45px' }}
                            id="menu-appbar"
                            anchorEl={anchorElUser}
                            anchorOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            open={Boolean(anchorElUser)}
                            onClose={handleCloseUserMenu}
                        >
                            {settings.map((setting) => (
                                <MenuItem key={setting} onClick={() => {handleCloseUserMenu(setting)}}>
                                    <Typography textAlign="center">{setting}</Typography>
                                </MenuItem>
                            ))}
                        </Menu>
                    </Box>
                </Toolbar>
            </Container>
        </AppBar>
            <Container>
                <Box sx={{ my: 3 }} textAlign = "left" >
                    <Card sx={{ minWidth: 275 }}>
                        <CardContent>
                            <Typography sx={{ fontSize: 20 }} color="info.main" gutterBottom>
                                Your setup
                            </Typography>
                            <Typography sx={{ mb: 1.5 }} color="text.secondary">
                                Step 1 - Download kubectl
                            </Typography>
                            <Stack direction="row" spacing={2}>
                                <Button variant="contained" href={`https://storage.googleapis.com/kubernetes-release/release/${kubernetes.serverVersion}/bin/darwin/amd64/kubectl`}>macOS</Button>
                                <Button variant="contained" href={`https://storage.googleapis.com/kubernetes-release/release/${kubernetes.serverVersion}/bin/windows/amd64/kubectl.exe`}>Windows</Button>
                                <Button variant="contained" href={`https://storage.googleapis.com/kubernetes-release/release/${kubernetes.serverVersion}/bin/linux/amd64/kubectl`}>Linux</Button>
                            </Stack>
                            <Typography sx={{ mb: 1.5, mt: 2 }} color="text.secondary">
                                Step 2 - Download kubeconfig
                            </Typography>
                                <code>
                                    #Using curl<br />
                                    {`curl -s ${window.location.protocol}//${window.location.host}/kubeconfig/${user.id} >> ~/.kube/config`}
                                </code>
                            <br />
                            <br />
                            <Button variant="contained" href={`${window.location.protocol}//${window.location.host}/kubeconfig/${user.id}`}>Download kubeconfig</Button>
                        </CardContent>

                    </Card>

                </Box>
                <Box sx={{ my: 3 }} textAlign = "left" >
                    <Card sx={{ minWidth: 275 }}>
                        <CardContent>
                            <Typography sx={{ fontSize: 20 }} color="info.main" gutterBottom>
                                Instructions
                            </Typography>

                            <Instructions />
                        </CardContent>

                    </Card>

                </Box>
            </Container>
        </Box>
    );
};
export default ResponsiveAppBar;
