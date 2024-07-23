# Appointment Creator Service

Service responsible to create the appointments

# Local Development

## Requirements

- [Kubernetes](https://kubernetes.io/)
- [AWS CLI](https://aws.amazon.com/cli/)

## Manual deployment

### Attention

Before deploying the service, make sure to set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.

Be aware that this process will take a few minutes (~4 minutes) to be completed.

To deploy the service manually, run the following commands in order:

```bash
make init
make check # this will execute fmt, validate and plan
make apply
```

To destroy the service, run the following command:

```bash
make destroy
```

## Automated deployment

The automated deployment is triggered by a GitHub Action.

# Consume Appointment Events

## CreateAppointment - Success flow
- [ ] Consume queue
- [ ] Check if the appointment is already created
- [ ] If not, create the appointment
- [ ] Update the event with success status
- [ ] Delete the message from the queue

## CreateAppointment - Error flow
- [ ] Consume queue
- [ ] Check if the appointment is already created
- [ ] If yes, update the event with error status
- [ ] Delete the message from the queue

## UpdateAppointment - Success flow
- [ ] Consume queue
- [ ] Check if the appointment can be re-scheduled
- [ ] If yes, update the appointment with the new date
- [ ] Update the event with success status
- [ ] Delete the message from the queue

## UpdateAppointment - Error flow
- [ ] Consume queue
- [ ] Check if the appointment can be re-scheduled
- [ ] If no, cancel the appointment
- [ ] Update the event with error status
- [ ] Delete the message from the queue

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.