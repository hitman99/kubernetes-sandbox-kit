import React from 'react'
import {Button, Form, Grid, Header, Image, Segment, Message} from 'semantic-ui-react'
import {register} from '../utils/api'
import Instructions from '../instructions'
import kubelogo from '../static/kubelogo.png'

class RegistrationForm extends React.Component {

  loadState() {
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

  validateEmail(email) {
    return this.emailRegex.test(email);
  }

  persistState(regData) {
    localStorage.setItem("registrationData", JSON.stringify(regData))
  }

  canRegister() {
    const {regData, registered} = this.state;
    return regData.user.email.length > 5 && this.validateEmail(regData.user.email) && !registered;
  }

  constructor(props) {
    super(props);
    let regData = this.loadState();
    this.state = {
      regData,
      registered: regData.user.id !== '',
      isError: false,
      errorMessage: "",
      isLoading: false
    };
    this.emailRegex = /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,5})+$/;
  }

  handleInputChange(which, ev) {
    let regData = {...this.state.regData};
    regData.user[which] = ev.value;
    this.setState({regData})
  }

  async submit() {
    this.setState({isLoading: true, isError: false, errorMessage: ""});
    let res = await register(this.state.regData);
    try {
      if (res.success !== true) {
        this.setState({isError: true, isLoading: false, errorMessage: res.message});
      } else {
        this.persistState(res);
        this.setState({regData: res, registered: true, isLoading: false, isError: false});
      }
    } catch (err) {
      this.setState({isError: true, isLoading: false, errorMessage: "something went wrong, please try again"});
    }
  }

  render() {
    const {user, kubernetes} = this.state.regData;
    const {registered, isError, isLoading, errorMessage} = this.state;
    let err;
    let formOrCard;
    if (isError) {
      err =
        <Message negative>
          <Message.Header>Sorry, there was an error</Message.Header>
          <p>{errorMessage}</p>
        </Message>;
    }
    if (!registered) {
      formOrCard =
        <React.Fragment>
          <Header as='h1' color='blue' textAlign='center'>
            <Image src={kubelogo}/> kubernetes sandbox
          </Header>
          <Form size='large'>
            <Segment>
              {err}
              <Form.Input
                inverted
                fluid icon='mail'
                iconPosition='left'
                placeholder='Email address'
                value={user.email}
                onChange={(e, d) => {
                  this.handleInputChange('email', d)
                }}
                disabled={registered}
                error={isError}
              />
              <Button loading={isLoading} color='blue' fluid size='large' onClick={() => {
                this.submit()
              }} disabled={!this.canRegister()}>
                Register
              </Button>
            </Segment>
          </Form>
        </React.Fragment>
    } else {
      formOrCard =
        <Instructions userData={this.state.regData}/>
    }
    let styles = {
      registered: {
        marginLeft: 50,
        marginRight: 50
      },
      unregistered: {
        maxWidth: 550
      }
    };
    return (
      <Grid textAlign='center' style={{height: '100vh'}} verticalAlign='middle'>
        <Grid.Column style={registered ? styles.registered : styles.unregistered}>
          {formOrCard}
        </Grid.Column>
      </Grid>
    )
  }
}

export default RegistrationForm