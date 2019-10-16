import React, { Component } from 'react';
import { func, object, bool } from 'prop-types';
import { upper } from 'change-case';
import { pipe, map, sort as rSort } from 'ramda';
import memoize from 'fast-memoize';
import TableRow from '@material-ui/core/TableRow';
import TableCell from '@material-ui/core/TableCell';
import Typography from '@material-ui/core/Typography';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import DeleteIcon from 'mdi-react/DeleteIcon';
import OpenInNewIcon from 'mdi-react/OpenInNewIcon';
// import { awsProvisionerWorkerType } from '../../utils/prop-types';
import sort from '../../utils/sort';
import DataTable from '../DataTable';
import TableCellItem from '../TableCellItem';
import StatusLabel from '../StatusLabel';
import DateDistance from '../DateDistance';

const awsConsoleUrl = 'https://console.aws.amazon.com/ec2/v2/home';
const sorted = pipe(
  rSort((a, b) => sort(a.id, b.id)),
  map(({ id }) => id)
);

export default class Ec2ResourcesTable extends Component {
  static defaultProps = {
    actionLoading: false,
    awsState: null,
  };

  static propTypes = {
    /** A GraphQL awsProvisionerWorkerType response. */
    // workerType: awsProvisionerWorkerType,
    /** A Graphql awsProvisionerWorkerTypeState response */
    awsState: object,
    /** Callback function fired when an EC2 instance is terminated. */
    onTerminateInstance: func.isRequired,
    /** If true, action buttons (e.g., terminate) will be disabled. */
    actionLoading: bool,
  };

  state = {
    sortBy: null,
    sortDirection: null,
  };

  createSortedInstances = memoize(
    (awsState, sortBy, sortDirection) => {
      if (!awsState) {
        return null;
      }

      return [...awsState.instances].sort((a, b) => {
        const firstElement = sortDirection === 'desc' ? b[sortBy] : a[sortBy];
        const secondElement = sortDirection === 'desc' ? a[sortBy] : b[sortBy];

        return sort(firstElement, secondElement);
      });
    },
    {
      serializer: ([awsState, sortBy, sortDirection]) => {
        if (!awsState) {
          return null;
        }

        const ids = sorted(awsState.instances);

        return `${ids.join('-')}-${sortBy}-${sortDirection}`;
      },
    }
  );

  handleHeaderClick = header => {
    const toggled = this.state.sortDirection === 'desc' ? 'asc' : 'desc';
    const sortDirection = this.state.sortBy === header.id ? toggled : 'desc';

    this.setState({ sortBy: header.id, sortDirection });
  };

  render() {
    const {
      onTerminateInstance,
      actionLoading,
      awsState /* , workerType */,
    } = this.props;
    const { sortBy, sortDirection } = this.state;
    const sortedInstances = this.createSortedInstances(
      awsState,
      sortBy,
      sortDirection
    );
    const iconSize = 16;
    const headers = [
      { label: 'Instance ID', id: 'id', type: 'string' },
      {
        label: 'State',
        id: 'state',
        type: 'string',
      },
      {
        label: 'Availability Zone',
        id: 'zone',
        type: 'string',
      },
      {
        label: 'Launch Time',
        id: 'launch',
        type: 'string',
      },
      {
        label: '',
        id: '',
        type: 'undefined',
      },
    ];

    return (
      <DataTable
        items={sortedInstances}
        headers={headers}
        sortByLabel={sortBy}
        sortDirection={sortDirection}
        onHeaderClick={this.handleHeaderClick}
        noItemsMessage="No EC2 resources"
        renderRow={instance => (
          <TableRow key={instance.id}>
            <TableCell>
              <TableCellItem
                button
                component="a"
                target="_blank"
                rel="noopener noreferrer"
                href={`${awsConsoleUrl}?region=${instance.region}#Images:visibility=owned-by-me;imageId=${instance.id};sort=name`}>
                <ListItemText primary={instance.id} />
                <OpenInNewIcon size={iconSize} />
              </TableCellItem>
            </TableCell>
            <TableCell>
              <StatusLabel state={upper(instance.state)} />
            </TableCell>
            <TableCell>
              <Typography>{instance.zone}</Typography>
            </TableCell>
            <TableCell>
              <DateDistance from={instance.launch} />
            </TableCell>
            <TableCell>
              <Tooltip placement="bottom" title="Terminate">
                <span>
                  <IconButton
                    disabled={actionLoading}
                    onClick={() => onTerminateInstance(instance)}>
                    <DeleteIcon size={18} />
                  </IconButton>
                </span>
              </Tooltip>
            </TableCell>
          </TableRow>
        )}
      />
    );
  }
}
