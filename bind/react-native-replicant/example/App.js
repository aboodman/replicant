/**
 * Sample React Native App
 *
 * adapted from App.js generated by the following command:
 *
 * react-native init example
 *
 * https://github.com/facebook/react-native
 */
 
import React, { Component } from 'react';
import { Platform, StyleSheet, Text, View, Button } from 'react-native';
import Replicant from './replicant.js';
import bundle from './bundle.js';

export default class App extends Component<{}> {
  state = {
    root: '',
  };
  async componentDidMount() {
    this._replicant = new Replicant('https://replicate.to/serve/react-native-test');
    await this._replicant.putBundle(bundle);

    const root = await this._replicant.root();
    this.setState({
      root,
    });

    await this._replicant.exec('addTodo', ['1', 'Pickup Abby', 1, false]);
    await this._replicant.exec('addTodo', ['1', 'Pickup Sam', 2, false]);
    console.warn('all todos', await this._replicant.exec('getAllTodos'));

    await this._replicant.exec('deleteAllTodos');
    console.warn('after delete', await this._replicant.exec('getAllTodos'));
  }
  async _handleSync() {
    const result = await Replicant.dispatch('sync', JSON.stringify({remote: 'https://replicate.to/serve/react-native-test'}));
    console.log('Sync result was', result);
  }
  render() {
    return (
      <View style={styles.container}>
        <Text style={styles.welcome}>☆Replicant example☆</Text>
        <Text style={styles.instructions}>Current root: {this.state.root}</Text>
        <Button onPress={this._handleSync} title="Sync"/>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F5FCFF',
  },
  welcome: {
    fontSize: 20,
    textAlign: 'center',
    margin: 10,
  },
  instructions: {
    textAlign: 'center',
    color: '#333333',
    marginBottom: 5,
  },
});
