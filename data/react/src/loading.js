import React from 'react';
import ReactLoading from 'react-loading';

const Loading = () => (
    <div align='center'>
	<ReactLoading type={'spin'} color={'#3498db'} height={'10%'} width={'10%'} />
    </div>
);

export default Loading;