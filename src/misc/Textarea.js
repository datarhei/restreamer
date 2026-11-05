import React, { useContext } from 'react';

import { useLingui } from '@lingui/react';
import { t } from '@lingui/macro';
import DownloadIcon from '@mui/icons-material/Download';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import IconButton from '@mui/material/IconButton';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import Stack from '@mui/material/Stack';

import NotifyContext from '../contexts/Notify';
import Palette from '../theme/base/palette';
import TextareaModal from './modals/Textarea';

export default function Component(props) {
	const { i18n } = useLingui();

	const [$modal, setModal] = React.useState(false);

	const notify = useContext(NotifyContext);
	const textAreaRef = React.createRef();

	const { content } = props;

	React.useEffect(() => {
		scrollTo();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [content]);

	const handleCopy = async () => {
		let success = false;

		if (!navigator.clipboard) {
			textAreaRef.current.select();

			try {
				success = document.execCommand('copy');
			} catch (err) {}

			textAreaRef.current.setSelectionRange(0, 0);
		} else {
			success = await writeText(navigator.clipboard.writeText(textAreaRef.current.value));
		}

		if (success === true) {
			notify.Dispatch('success', 'clipboard', i18n._(t`Data copied to clipboard`));
		} else {
			notify.Dispatch('error', 'clipboard', i18n._(t`Error while copying data to clipboard`));
		}
	};

	const handleModal = () => {
		setModal(!$modal);
	};

	const writeText = (promise) => {
		return promise
			.then(() => true)
			.catch((err) => {
				console.warn(err);
				return false;
			});
	};

	const scrollTo = () => {
		if (props.scrollTo === 'bottom') {
			textAreaRef.current.scrollTop = textAreaRef.current.scrollHeight;
		}

		return;
	};

	const handleDownload = () => {
		var element = document.createElement('a');
		element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(props.value));
		element.setAttribute('download', props.downloadName);

		element.style.display = 'none';
		document.body.appendChild(element);

		element.click();

		document.body.removeChild(element);
	};

	let allowCopy = props.allowCopy;
	if (props.value.length === 0 || props.disabled === true) {
		allowCopy = false;
	}

	let allowModal = props.allowModal;
	if (props.value.length === 0 || props.disabled === true) {
		allowModal = false;
	}

	let allowDownload = props.allowDownload;
	if (props.value.length === 0 || props.disabled === true || props.downloadName.length === 0) {
		allowDownload = false;
	}

	let textAreaDivStyle = {
		width: '100%',
	};
	let textAreaStyle = {
		lineHeight: 1.3,
		fontSize: 14,
	};
	let actionButton = {
		float: 'right',
		padding: '0.5em',
		backgroundColor: Palette.background.footer1,
		'&:hover': {
			backgroundColor: Palette.background.footer2,
		},
	};
	if (props.rows === 1) {
		textAreaStyle = {
			...textAreaStyle,
			height: 18 * props.rows + 9.5 + 'px',
			overflowY: 'hidden',
			marginBottom: '0em',
			marginTop: '0em',
		};
		textAreaDivStyle = {
			...textAreaDivStyle,
			padding: '1em 0em 1em 0em',
		};
		actionButton = {
			...actionButton,
			bottom: '-1.8em',
			marginRight: '.3em',
			marginTop: '-2em',
		};
	} else {
		textAreaStyle = {
			...textAreaStyle,
			height: 18 * props.rows + 8 + 'px',
		};
		textAreaDivStyle = {
			...textAreaDivStyle,
		};
		actionButton = {
			...actionButton,
			bottom: '-3em',
			marginRight: '.8em',
			marginTop: '-2em',
		};
	}

	return (
		<React.Fragment>
			<Stack direction="column" justifyContent="center" alignItems="flex-start" spacing={1} style={textAreaDivStyle}>
				<Stack direction="column" justifyContent="flex-start" alignItems="flex-end" spacing={0} width="100%">
					{allowCopy && (
						<IconButton size="small" onClick={handleCopy} style={actionButton}>
							<FileCopyIcon fontSize="small" />
						</IconButton>
					)}
					{allowDownload && (
						<IconButton size="small" onClick={handleDownload} style={actionButton}>
							<DownloadIcon fontSize="small" />
						</IconButton>
					)}
					{allowModal && (
						<IconButton size="small" onClick={handleModal} style={actionButton}>
							<OpenInNewIcon fontSize="small" />
						</IconButton>
					)}
				</Stack>
				<textarea
					style={textAreaStyle}
					ref={textAreaRef}
					rows={props.rows}
					value={props.value}
					readOnly={props.readOnly}
					disabled={props.disabled}
					onChange={props.onChange}
				/>
			</Stack>
			<TextareaModal
				open={$modal}
				onClose={handleModal}
				title={props.title}
				onHelp={props.onHelp}
				rows={props.rows}
				value={props.value}
				readOnly={props.readOnly}
				disabled={props.disabled}
				onChange={props.onChange}
				scrollTo={props.scrollTo}
				{...props}
				allowModal={false}
			/>
		</React.Fragment>
	);
}

Component.defaultProps = {
	title: '',
	rows: 20,
	value: '',
	readOnly: true,
	allowCopy: true,
	allowModal: false,
	allowDownload: false,
	downloadName: '',
	disabled: false,
	scrollTo: 'top',
	onChange: function (value) {},
	onHelp: null,
};
