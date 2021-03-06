import React from 'react';
import classNames from 'classnames';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import Table, {
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
  TableSortLabel,
} from 'material-ui/Table';
import Typography from 'material-ui/Typography';
import Paper from 'material-ui/Paper';
import Toolbar from 'material-ui/Toolbar';
import Button from 'material-ui/Button';
import Tooltip from 'material-ui/Tooltip';
import FileUpload from 'material-ui-icons/FileUpload';
import FileDownload from 'material-ui-icons/FileDownload';
import FilterListIcon from 'material-ui-icons/FilterList';
import { lighten } from 'material-ui/styles/colorManipulator';
import Loading from 'app/components/Loading';
import ImportDialog from 'app/components/dialogs/Import';
import withAuth from 'app/components/withAuth';
import AuthService from 'app/components/AuthService';

const columnData = [
  { id: 'date', numeric: false, disablePadding: true, label: 'Date' },
  { id: 'type', numeric: false, disablePadding: false, label: 'Type' },
  { id: 'currency_pair', numeric: false, disablePadding: false, label: 'Currency' },
  { id: 'network', numeric: false, disablePadding: false, label: 'Network' },
  { id: 'quantity', numeric: true, disablePadding: false, label: 'Quantity' },
  { id: 'price', numeric: true, disablePadding: false, label: 'Price' },
  { id: 'fee', numeric: true, disablePadding: false, label: 'Fee' },
  { id: 'total', numeric: true, disablePadding: false, label: 'Total' }
];

class OrderHistoryHead extends React.Component {

  createSortHandler = property => event => {
    this.props.onRequestSort(event, property);
  };

  render() {
    const { order, orderBy, numSelected, rowCount } = this.props;

    return (
      <TableHead>
        <TableRow>
          {columnData.map(column => {
            return (
              <TableCell
                key={column.id}
                numeric={column.numeric}
                padding={column.disablePadding ? 'none' : 'default'}
                sortDirection={orderBy === column.id ? order : false}>
                <Tooltip
                  title="Sort"
                  placement={column.numeric ? 'bottom-end' : 'bottom-start'}
                  enterDelay={300}>
                  <TableSortLabel
                    active={orderBy === column.id}
                    direction={order}
                    onClick={this.createSortHandler(column.id)}>
                    {column.label}
                  </TableSortLabel>
                </Tooltip>
              </TableCell>
            );
          }, this)}
        </TableRow>
      </TableHead>
    );
  }
}

OrderHistoryHead.propTypes = {
  numSelected: PropTypes.number.isRequired,
  onRequestSort: PropTypes.func.isRequired,
  order: PropTypes.string.isRequired,
  orderBy: PropTypes.string.isRequired,
  rowCount: PropTypes.number.isRequired,
};

const styles = theme => ({
  root: {
    flex: 1,
    flexGrow: 1,
    paddingLeft: '1%',
    width: '99%',
    marginTop: '68px'
    //marginTop: theme.spacing.unit * 8,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
  table: {
    width: '100%'
  },
  tableWrapper: {
    overflowX: 'auto',
  },
  currencyIcon: {
    paddingLeft: '5px',
    width: '16px',
    height: '16px',
    float: 'right'
  },
  loadingContainer: {
    marginTop: '68px'
  },
  leftIcon: {
    marginRight: theme.spacing.unit,
    color: 'grey'
  },
  buttonText: {
    color: 'grey',
    fontSize: '12px'
  }
});

class OrderHistory extends React.Component {

  constructor(props, context) {
    super(props, context);
    this.Auth = new AuthService();
    this.state = {
      loading: true,
      importDialog: false,
      local_currency: this.Auth.getUser().local_currency,
      order: 'asc',
      orderBy: 'date',
      selected: [],
      data: [],
      page: 0,
      rowsPerPage: 10
    }
    this.fetchOrderHistory = this.fetchOrderHistory.bind(this)
  }

  handleRequestSort = (event, property) => {
    const orderBy = property;
    let order = 'desc';

    if (this.state.orderBy === property && this.state.order === 'desc') {
      order = 'asc';
    }

    const data =
      order === 'desc'
        ? this.state.data.sort((a, b) => (b[orderBy] < a[orderBy] ? -1 : 1))
        : this.state.data.sort((a, b) => (a[orderBy] < b[orderBy] ? -1 : 1));

    this.setState({ data, order, orderBy });
  };

  handleChangePage = (event, page) => {
    this.setState({ page });
  };

  handleChangeRowsPerPage = event => {
    this.setState({ rowsPerPage: event.target.value });
  };

  isSelected = id => this.state.selected.indexOf(id) !== -1;

  componentDidMount() {
    this.fetchOrderHistory()
	}

  fetchOrderHistory() {
    this.Auth.fetchOrderHistory()
      .then(function (response) {
        console.log(response);
        if(response.success) {
  		    this.setState({
            loading: false,
            data: response.payload
          })
        }
      }.bind(this))
  }

  currencyIcon(currency) {
    return "images/crypto/128/" + currency.toLowerCase() + ".png";
  }

  render() {
    const { classes } = this.props;
    const { data, order, orderBy, selected, rowsPerPage, page } = this.state;
    const emptyRows = rowsPerPage - Math.min(rowsPerPage, data.length - page * rowsPerPage);

    return (
      <Paper className={classes.root}>
        {this.state.loading &&
          <div className={classes.loadingContainer}>
            <Loading text="Loading order history..." />
          </div>
        }
        <div className={classes.tableWrapper}>
          <Table className={classes.table}>
            <OrderHistoryHead
              numSelected={selected.length}
              order={order}
              orderBy={orderBy}
              onSelectAllClick={this.handleSelectAllClick}
              onRequestSort={this.handleRequestSort}
              rowCount={data.length}
            />
            <TableBody className={classes.tableBody}>
              {data.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map(n => {
                const isSelected = this.isSelected(n.id);

                const columnData = [
                  { id: 'date', numeric: false, disablePadding: false, label: 'Date' },
                  { id: 'type', numeric: false, disablePadding: false, label: 'Type' },
                  { id: 'currency', numeric: false, disablePadding: false, label: 'Currency' },
                  { id: 'network', numeric: false, disablePadding: false, label: 'Network' },
                  { id: 'quantity', numeric: true, disablePadding: false, label: 'Quantity' },
                  { id: 'price', numeric: true, disablePadding: false, label: 'Price' },
                  { id: 'fee', numeric: true, disablePadding: false, label: 'Fee' },
                  { id: 'total', numeric: true, disablePadding: false, label: 'Total' },
                ];

                return (
                  <TableRow key={n.id}>
                    <TableCell padding="none">{new Date(n.date).customFormat()}</TableCell>
                    <TableCell>{n.type}</TableCell>
                    <TableCell>{n.currency_pair.base}-{n.currency_pair.quote}</TableCell>
                    <TableCell>{n.network_display_name}</TableCell>
                    <TableCell numeric>{n.quantity}
                      <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.quantity_currency)}
                           title={n.quantity_currency} />
                    </TableCell>
                    <TableCell numeric>{n.price.formatCurrency(n.price_currency)}
                      <img className={classes.currencyIcon}
                         src={this.currencyIcon(n.price_currency)}
                         title={n.price_currency} />
                    </TableCell>
                    <TableCell numeric>{n.fee.formatCurrency(n.fee_currency)}
                      <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.fee_currency)}
                           title={n.fee_currency} />
                    </TableCell>
                    <TableCell numeric>{n.total.formatCurrency(n.total_currency)}
                      <img className={classes.currencyIcon}
                           src={this.currencyIcon(n.total_currency)}
                           title={n.total_currency} />
                    </TableCell>
                  </TableRow>
                );
              })}
              {emptyRows > 0 && (
                <TableRow style={{ height: 49 * emptyRows }}>
                  <TableCell colSpan={8} />
                </TableRow>
              )}
            </TableBody>
            <TableFooter>
              <TableRow>
                <TablePagination
                  colSpan={8}
                  count={data.length}
                  rowsPerPage={rowsPerPage}
                  page={page}
                  backIconButtonProps={{
                    'aria-label': 'Previous Page',
                  }}
                  nextIconButtonProps={{
                    'aria-label': 'Next Page',
                  }}
                  onChangePage={this.handleChangePage}
                  onChangeRowsPerPage={this.handleChangeRowsPerPage}
                />
              </TableRow>
            </TableFooter>
          </Table>
        </div>
      </Paper>
    );
  }
}

OrderHistory.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withAuth(withStyles(styles)(OrderHistory));
